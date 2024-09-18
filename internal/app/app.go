package app

import (
	"fmt"
	"github.com/daemondxx/lks_back/entity"
	authpb "github.com/daemondxx/lks_back/gen/pb/go/auth"
	orderpb "github.com/daemondxx/lks_back/gen/pb/go/order"
	userpb "github.com/daemondxx/lks_back/gen/pb/go/user"
	"github.com/daemondxx/lks_back/internal/api/lks"
	"github.com/daemondxx/lks_back/internal/dao"
	"github.com/daemondxx/lks_back/internal/servers"
	"github.com/daemondxx/lks_back/internal/services/authchecker"
	"github.com/daemondxx/lks_back/internal/services/collector"
	"github.com/daemondxx/lks_back/internal/services/notifier"
	"github.com/daemondxx/lks_back/internal/services/notifier/implementation"
	"github.com/daemondxx/lks_back/internal/services/order"
	"github.com/daemondxx/lks_back/internal/services/user"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type LKSApp struct {
	db       *gorm.DB
	log      *zerolog.Logger
	listener net.Listener
	server   *grpc.Server
}

func NewApp(cfg Config, log *zerolog.Logger) (*LKSApp, error) {
	app := &LKSApp{
		log: log,
	}

	if err := app.initDB(cfg.Database); err != nil {
		return nil, fmt.Errorf("init database error: %e", err)
	}

	if err := app.initServices(cfg); err != nil {
		return nil, fmt.Errorf("init services error: %e", err)
	}

	return app, nil
}

func (a *LKSApp) initDB(cfg DBConfig) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connect to db error: %e", err)
	}

	a.db = db

	if err := db.AutoMigrate(&entity.User{}); err != nil {
		return fmt.Errorf("migrate user entity error: %e", err)
	}

	if err := db.AutoMigrate(&entity.FlightInfo{}); err != nil {
		return fmt.Errorf("migrate flightinfo entity error: %e", err)
	}

	if err := db.AutoMigrate(&entity.Order{}); err != nil {
		return fmt.Errorf("migrate order entity error: %e", err)
	}

	if err := db.AutoMigrate(&entity.OrderItem{}); err != nil {
		return fmt.Errorf("migrate order item entity error: %e", err)
	}

	if err := db.AutoMigrate(&entity.Flight{}); err != nil {
		return fmt.Errorf("migrate flight entity error: %e", err)
	}

	return nil
}

func (a *LKSApp) initServices(cfg Config) error {
	userDAO := dao.NewUserDAO(a.db)

	uServLog := a.log.With().Str("service", "user_service").Logger()
	userServ := user.NewUserService(userDAO, &uServLog)

	//fInfoDao := dao.NewFlightInfoDAO(a.db)
	//fInfoApiLog := a.log.With().Str("service", "flightaware_service").Logger()
	//fInfoApi, err := flightaware.NewFlightInfoAPI(&flightaware.ApiConfig{
	//	MaxTabCount: cfg.FlightAPI.MaxTabCount,
	//	Debug: cfg.FlightAPI.DebugMode,
	//}, &fInfoApiLog)
	//if err != nil {
	//	return fmt.Errorf("create flightinfo api error: %e", err)
	//}
	//fInfoServLog := a.log.With().Str("service", "flight_info_service").Logger()
	//fInfoServ := flight.NewFlightInfoService(fInfoDao, fInfoApi, &fInfoServLog)

	lksApiCachedLog := a.log.With().Str("service", "lks_api_cached").Logger()
	oDao := dao.NewOrderDAO(a.db)
	lksApiCached := lks.NewLksAPI(&lks.LksAPIConfig{
		AuthWorkerPoolSize: uint(cfg.LKSApi.WorkerPoolSize),
		Debug:              cfg.LKSApi.DebugMode,
	}, lks.NewSimpleCookieCache(), &lksApiCachedLog)

	oServLog := a.log.With().Str("service", "order_service").Logger()
	oServ := order.NewOrderService(oDao, lksApiCached, userServ, &oServLog)

	lksApiLog := a.log.With().Str("service", "lks_api_for_checker").Logger()
	lksApi := lks.NewLksAPI(&lks.LksAPIConfig{
		AuthWorkerPoolSize: 1,
		Debug:              false,
	}, nil, &lksApiLog)
	checker := authchecker.NewAuthCheckerClient(lksApi)

	if err := a.checkKafkaConnection(cfg.Notifier.Addr); err != nil {
		return err
	}

	kWriter := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Notifier.Addr...),
		Topic:                  cfg.Notifier.Topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: false,
	}

	notifyLog := a.log.With().Str("service", "notifier").Logger()
	notifyServ := notifier.NewNotifierService(kWriter, &notifyLog)

	collectorNotify := implementation.NewCollectorNotifier(notifyServ)

	collectorLog := a.log.With().Str("service", "collector").Logger()
	collectorServ := collector.NewCollectorService(userDAO, oServ, collectorNotify, &collector.Config{
		MaxAttempts:     6 * 3,
		MinTimeoutRetry: 10 * time.Minute,
	}, &collectorLog)

	autoCollectLog := a.log.With().Str("service", "auto_collector").Logger()
	autoCollectServ, err := collector.NewAutoWorker(collectorServ, &collector.AutoWorkerConfig{
		ActualOrderCronList: cfg.AutoCollector.ActualOrderCronList,
		MonthOrderCronList:  nil,
		TaskTimeout:         cfg.AutoCollector.TaskTimeout,
	}, &autoCollectLog)

	autoCollectServ.Start()

	a.server = grpc.NewServer()

	authServerLog := a.log.With().Str("service", "auth_server").Logger()
	authServer := servers.NewAuthServer(userServ, checker, &authServerLog)
	authpb.RegisterAuthServiceServer(a.server, authServer)

	userServerLog := a.log.With().Str("service", "user_server").Logger()
	userServer := servers.NewUserServer(userServ, &userServerLog)
	userpb.RegisterUserServiceServer(a.server, userServer)

	orderServer := servers.NewOrderServer(oServ)
	orderpb.RegisterOrderServiceServer(a.server, orderServer)

	reflection.Register(a.server)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPC.Port))
	if err != nil {
		return fmt.Errorf("listen tcp error: %e", err)
	}

	a.listener = listener

	return nil
}

func (a *LKSApp) checkKafkaConnection(addresses []string) error {
	for _, address := range addresses {
		conn, err := kafka.Dial("tcp", address)
		if err != nil {
			a.log.Warn().Msg(fmt.Sprintf("failed connection to kafka %s: %e", address, err))
		} else {
			conn.Close()
			return nil
		}
	}
	return fmt.Errorf("connection to kafka failed")
}

func (a *LKSApp) Run() error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	go func() {
		<-sigCh
		a.server.Stop()
		a.listener.Close()
	}()

	if err := a.server.Serve(a.listener); err != nil {
		return fmt.Errorf("start grpc server error: %e", err)
	} else {
		return nil
	}
}

func (a *LKSApp) Stop() error {
	a.server.Stop()
	return a.listener.Close()
}
