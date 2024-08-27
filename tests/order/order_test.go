package order

import (
	"context"
	"fmt"
	"github.com/daemondxx/lks_back/entity"
	"github.com/daemondxx/lks_back/internal/api/lks"
	"github.com/daemondxx/lks_back/internal/dao"
	"github.com/daemondxx/lks_back/internal/services/order"
	"github.com/daemondxx/lks_back/internal/services/user"
	mock_order "github.com/daemondxx/lks_back/tests/mocks"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"strings"
	"testing"
	"time"
)

type dbHelper struct {
	db *gorm.DB
}

func (h *dbHelper) getCountOrders() (int, error) {
	var orders []entity.Order
	if err := h.db.Where("id > 0").Find(&orders).Error; err != nil {
		return 0, err
	} else {
		return len(orders), nil
	}
}

func (h *dbHelper) getCountOrderItems() (int, error) {
	var items []entity.OrderItem
	if err := h.db.Where("id > 0").Find(&items).Error; err != nil {
		return 0, err
	} else {
		return len(items), nil
	}
}

func (h *dbHelper) getCountFlights() (int, error) {
	var flights []entity.Flight
	if err := h.db.Where("id > 0").Find(&flights).Error; err != nil {
		return 0, err
	} else {
		return len(flights), nil
	}
}

type OrderIntegrationSuite struct {
	suite.Suite
	userID uint
	db     *gorm.DB
	log    *zerolog.Logger
	lks    *mock_order.MockLksAPI
	serv   *order.Service
	h      *dbHelper
}

func (s *OrderIntegrationSuite) SetupSuite() {

	log := zerolog.New(os.Stdin).Level(zerolog.ErrorLevel)
	s.log = &log

	if err := godotenv.Load("../.env"); err != nil {
		s.T().Fatal(fmt.Sprintf("read config file error: %e", err))
	}

	err := s.initDB()
	if err != nil {
		s.T().Fatal(fmt.Sprintf("db init error: %e", err))
	}
	if err := s.clearDB(); err != nil {
		s.T().Fatal(fmt.Sprintf("clear db setup suite error: %e", err))
	}

	s.h = &dbHelper{db: s.db}

	uDAO := dao.NewUserDAO(s.db)
	uServ := user.NewUserService(uDAO, s.log)
	oDAO := dao.NewOrderDAO(s.db)
	l := mock_order.NewMockLksAPI(s.T())
	s.lks = l
	s.serv = order.NewOrderService(oDAO, l, uServ, s.log)
}

func (s *OrderIntegrationSuite) initDB() error {
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	usr := os.Getenv("PG_USER")
	pass := os.Getenv("PG_PASSWORD")
	name := os.Getenv("PG_DB")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, usr, pass, name, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("open connection error: %e", err)
	}

	if err := db.AutoMigrate(&entity.User{}); err != nil {
		return fmt.Errorf("migrate user entity error: %e", err)
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

	s.db = db

	return nil
}

func (s *OrderIntegrationSuite) clearDB() error {
	if err := s.db.Unscoped().Where("id > 0").Delete(&entity.User{}).Error; err != nil {
		return err
	}
	if err := s.db.Unscoped().Where("id > 0").Delete(&entity.Flight{}).Error; err != nil {
		return err
	}
	if err := s.db.Unscoped().Where("id > 0").Delete(&entity.OrderItem{}).Error; err != nil {
		return err
	}
	if err := s.db.Unscoped().Where("id > 0").Delete(&entity.Order{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *OrderIntegrationSuite) BeforeTest(sName string, tName string) {
	u := &entity.User{
		AccordLogin:    "test",
		LKSLogin:       "test",
		AccordPassword: "test",
		LKSPassword:    "test",
	}
	if err := s.db.Create(u).Error; err != nil {
		s.T().Fatal(fmt.Sprintf("[%s] - create new user error: %e", tName, err))
	}
	s.userID = u.ID
}

func (s *OrderIntegrationSuite) AfterTest(sName string, tName string) {
	if err := s.clearDB(); err != nil {
		s.T().Fatal(err)
	}

}

func (s *OrderIntegrationSuite) TestGetFirstOrder() {
	expectItems := getFirstDuty(new(time.Time))
	fn := s.lks.EXPECT().GetActualDuty(mock.Anything, mock.Anything).Return(expectItems, nil)
	o, err := s.serv.GetActualOrder(context.Background(), s.userID)
	require.NoError(s.T(), err)

	require.Equal(s.T(), len(expectItems), len(o.Items))
	for i, d := range expectItems {
		assert.Equal(s.T(), len(strings.Split(d.FlightNumber, " ")), len(o.Items[i].Flights))
	}

	s.lks.AssertExpectations(s.T())
	fn.Unset()
}

func (s *OrderIntegrationSuite) TestGetCreatedOrder() {
	expectItems := getFirstDuty(new(time.Time))
	fn := s.lks.EXPECT().GetActualDuty(mock.Anything, mock.Anything).Return(expectItems, nil)
	defer func() {
		s.lks.AssertExpectations(s.T())
		fn.Unset()
	}()

	o1, err := s.serv.GetActualOrder(context.Background(), s.userID)
	require.NoError(s.T(), err)

	o2, err := s.serv.GetActualOrder(context.Background(), s.userID)
	require.NoError(s.T(), err)

	equalOrders(s.T(), o1, o2)

	countFlights := 0
	for _, d := range expectItems {
		countFlights += len(strings.Split(d.FlightNumber, " "))
	}
	s.assertsCountDBRows(1, len(expectItems), countFlights)
}

func (s *OrderIntegrationSuite) assertsCountDBRows(oCount int, iCount int, fCount int) {
	//should be create one order
	i, err := s.h.getCountOrders()
	require.NoError(s.T(), err)
	assert.Equal(s.T(), oCount, i)

	//should dont create duplicate order items
	i, err = s.h.getCountOrderItems()
	require.NoError(s.T(), err)
	assert.Equal(s.T(), iCount, i)

	//should dont create duplicate flights
	i, err = s.h.getCountFlights()
	require.NoError(s.T(), err)
	assert.Equal(s.T(), fCount, i)
}

func (s *OrderIntegrationSuite) TestGetOrderAfterCompletedFlight() {
	orderItems := getFirstDuty(&time.Time{})

	fn := s.lks.EXPECT().GetActualDuty(mock.Anything, mock.Anything).Return(orderItems, nil)

	o, err := s.serv.GetActualOrder(context.Background(), s.userID)
	require.NoError(s.T(), err)

	itemsAfterFlight := getFirstDuty(&time.Time{})

	s.lks.AssertExpectations(s.T())
	fn.Unset()
	fn = s.lks.EXPECT().GetActualDuty(mock.Anything, mock.Anything).Return(itemsAfterFlight[1:], nil)
	defer func() {
		s.lks.AssertExpectations(s.T())
		fn.Unset()
	}()

	oAfterFlight, err := s.serv.GetActualOrder(context.Background(), s.userID)
	require.NoError(s.T(), err)

	equalOrders(s.T(), o, oAfterFlight)
	countFlights := 0
	for _, d := range orderItems {
		countFlights += len(strings.Split(d.FlightNumber, " "))
	}
	s.assertsCountDBRows(1, len(orderItems), countFlights)
}

func (s *OrderIntegrationSuite) TestGetNewOrderWithOldOrderInDB() {
	items1 := getFirstDuty(&time.Time{})

	fn := s.lks.EXPECT().GetActualDuty(mock.Anything, mock.Anything).Return(items1, nil)
	o1, err := s.serv.GetActualOrder(context.Background(), s.userID)
	require.NoError(s.T(), err)
	s.lks.AssertExpectations(s.T())
	fn.Unset()

	items2 := getSecondDuty(&time.Time{})
	sumItems := make([]lks.CurrentDuty, 0, len(items1)+len(items2))
	sumItems = append(sumItems, items1...)
	sumItems = append(sumItems, items2...)

	fn = s.lks.EXPECT().GetActualDuty(mock.Anything, mock.Anything).Return(sumItems, nil)
	o2, err := s.serv.GetActualOrder(context.Background(), s.userID)
	require.NoError(s.T(), err)
	s.lks.AssertExpectations(s.T())
	fn.Unset()

	assert.NotEqual(s.T(), o1.ID, o2.ID)
	assert.Equal(s.T(), len(items1), len(o1.Items))
	assert.Equal(s.T(), len(items2), len(o2.Items))

	dutyCountFlights := 0
	for _, d := range items1 {
		dutyCountFlights += len(strings.Split(d.FlightNumber, " "))
	}

	oCountFlights1 := 0
	for _, i := range o1.Items {
		oCountFlights1 += len(i.Flights)
	}
	assert.Equal(s.T(), dutyCountFlights, oCountFlights1)

	dutyCountFlights = 0
	for _, d := range items2 {
		dutyCountFlights += len(strings.Split(d.FlightNumber, " "))
	}

	oCountFlights2 := 0
	for _, i := range o2.Items {
		oCountFlights2 += len(i.Flights)
	}
	assert.Equal(s.T(), dutyCountFlights, oCountFlights2)

	s.assertsCountDBRows(2, len(items1)+len(items2), oCountFlights1+oCountFlights2)
}

func (s *OrderIntegrationSuite) TestUpdateStatusOrderOnConfirmed() {
	itemsNotConfirm := getFirstDuty(nil)

	fn := s.lks.EXPECT().GetActualDuty(mock.Anything, mock.Anything).Return(itemsNotConfirm, nil)
	oNotConfirm, err := s.serv.GetActualOrder(context.Background(), s.userID)
	require.NoError(s.T(), err)
	s.lks.AssertExpectations(s.T())
	fn.Unset()

	assert.Equal(s.T(), entity.AwaitConfirmation, oNotConfirm.Status)

	confirmTime := time.Date(2024, 8, 8, 16, 15, 0, 0, &time.Location{})
	itemsConfirm := getFirstDuty(&confirmTime)
	fn = s.lks.EXPECT().GetActualDuty(mock.Anything, mock.Anything).Return(itemsConfirm, nil)
	oConfirm, err := s.serv.GetActualOrder(context.Background(), s.userID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), entity.Confirm, oConfirm.Status)

	var o []entity.Order
	err = s.db.Where("id > 0 AND status = 1").Find(&o).Error
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 1, len(o))
}

func (s *OrderIntegrationSuite) TestUpdateStatusOrderOnLimited() {
	item := getFirstDuty(&time.Time{})[:1]

	fn := s.lks.EXPECT().GetActualDuty(mock.Anything, mock.Anything).Return(item, nil)
	oConfirm, err := s.serv.GetActualOrder(context.Background(), s.userID)
	require.NoError(s.T(), err)
	s.lks.AssertExpectations(s.T())
	fn.Unset()

	assert.Equal(s.T(), entity.Confirm, oConfirm.Status)

	item[0].BlockDate = &time.Time{}
	fn = s.lks.EXPECT().GetActualDuty(mock.Anything, mock.Anything).Return(item, nil)
	oLimited, err := s.serv.GetActualOrder(context.Background(), s.userID)
	require.NoError(s.T(), err)
	s.lks.AssertExpectations(s.T())
	fn.Unset()

	assert.Equal(s.T(), entity.Limited, oLimited.Status)

	var o []entity.Order
	err = s.db.Where("id > 0 AND status = 2").Find(&o).Error
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 1, len(o))
}

func (s *OrderIntegrationSuite) TestGetEmptyOrder() {
	emptyItems := make([]lks.CurrentDuty, 0)
	fn := s.lks.EXPECT().GetActualDuty(mock.Anything, mock.Anything).Return(emptyItems, nil)

	o, err := s.serv.GetActualOrder(context.Background(), s.userID)
	require.NoError(s.T(), err)
	s.lks.AssertExpectations(s.T())
	fn.Unset()

	assert.Zero(s.T(), o.ID)
	assert.Equal(s.T(), 0, len(o.Items))

	s.assertsCountDBRows(0, 0, 0)
}

func (s *OrderIntegrationSuite) TestGetOrderWithoutFlights() {
	items := getNoFlightsDuty(&time.Time{})

	fn := s.lks.EXPECT().GetActualDuty(mock.Anything, mock.Anything).Return(items, nil)
	oNoFlights, err := s.serv.GetActualOrder(context.Background(), s.userID)
	require.NoError(s.T(), err)
	s.lks.AssertExpectations(s.T())
	fn.Unset()

	require.Equal(s.T(), len(items), len(oNoFlights.Items))
	for _, i := range oNoFlights.Items {
		require.Equal(s.T(), 1, len(i.Flights))
		assert.NotContains(s.T(), "AFL", i.Flights[0].FlightNumber)
	}
}

func TestIntegrationOrderSuite(t *testing.T) {
	suite.Run(t, new(OrderIntegrationSuite))
}

func equalOrders(t *testing.T, o1 entity.Order, o2 entity.Order) {
	assert.Equal(t, o1.ID, o2.ID)
	assert.Equal(t, o1.Status, o2.Status)
	assert.Equal(t, o1.UserID, o2.UserID)
	assert.Equal(t, len(o1.Items), len(o2.Items))
	for i, oItem := range o1.Items {
		assert.Equal(t, len(oItem.Flights), len(o2.Items[i].Flights))
	}
}
