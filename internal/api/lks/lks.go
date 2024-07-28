package lks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
	"time"
)

const (
	currentOrderUrl     = "https://lks.aeroflot.ru/AkkordOffice/GetCommitedRosters"
	perspectiveOrderUrl = "https://lks.aeroflot.ru/AkkordOffice/GetPerspectivePlan"
	archiveOrderUrl     = "https://lks.aeroflot.ru/AkkordOffice/GetAchievedDuties"
)

var (
	ErrAccordAuth    = errors.New("accord login or password incorrect")
	ErrLKSAuth       = errors.New("lks login or password incorrect")
	ErrExpiredCookie = errors.New("cookies is expired")
)

type AuthPayload struct {
	AccordLogin    string
	AccordPassword string
	LksLogin       string
	LksPassword    string
}

type LksAPIConfig struct {
	AuthWorkerPoolSize uint
	Debug              bool
}

type LksAPI struct {
	log   *zerolog.Logger
	cache CookieCache
	pool  *AuthWorkerPool
}

func NewLksAPI(cfg *LksAPIConfig, cache CookieCache, log *zerolog.Logger) *LksAPI {
	l := log.With().Str("service", "lks_api").Logger()
	pool := NewAuthWorkerPool(&AuthWorkerPoolConfig{
		PoolSize: cfg.AuthWorkerPoolSize,
		Debug:    cfg.Debug,
	}, log)
	lks := &LksAPI{
		log:   &l,
		cache: cache,
		pool:  pool,
	}
	return lks
}

func (a *LksAPI) GetActualDuty(ctx context.Context, p AuthPayload) ([]CurrentDuty, error) {
	log := a.getLogger("get_actual_duty")
	cookie, err := a.auth(ctx, p.AccordLogin, p.AccordPassword, p.LksLogin, p.LksPassword)
	if err != nil {
		return nil, err
	}

	log.Debug().Msg("prepare request")
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetRequestURI(currentOrderUrl)
	req.Header.SetMethod(fasthttp.MethodPost)
	for name, value := range cookie {
		req.Header.SetCookie(name, value)
	}
	req.Header.SetContentType("application/json; charset=UTF-8")
	req.SetBodyString(fmt.Sprintf("{\"staffNumber\":%s}", p.LksLogin))

	log.Debug().Msg("send request...")
	if err := fasthttp.DoTimeout(req, res, 20*time.Second); err != nil {
		log.Debug().Msg(fmt.Sprintf("unauthorized error: %e", ErrExpiredCookie))
		return nil, err
	}

	if status := res.StatusCode(); status != 200 {
		if status == 302 {
			log.Debug().Msg(fmt.Sprintf("unauthorized error: %e", ErrExpiredCookie))
			return nil, ErrExpiredCookie
		}
		log.Debug().Msg(fmt.Sprintf("send request returs status %d", err))
		return nil, fmt.Errorf("send request return status %d", status)
	}

	var resBody currentDutyResponse
	if err := json.Unmarshal(res.Body(), &resBody); err != nil {
		log.Debug().Msg(fmt.Sprintf("response body unmarshal error: %e", err))
		return nil, err
	}

	return resBody.Model.Duties, nil

}

func (a *LksAPI) GetPerspectiveDuty(p AuthPayload, month int, year int) ([]PerspectiveDuty, error) {
	log := a.getLogger("get_perspective_duty")

	cookie, err := a.auth(context.Background(), p.AccordLogin, p.AccordPassword, p.LksLogin, p.LksPassword)
	if err != nil {
		return nil, err
	}

	log.Debug().Msg("prepare request")
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetRequestURI(perspectiveOrderUrl)
	req.Header.SetMethod(fasthttp.MethodPost)
	for name, value := range cookie {
		req.Header.SetCookie(name, value)
	}
	req.Header.SetContentType("application/json; charset=UTF-8")

	var payload struct {
		Start    time.Time `json:"dateFrom"`
		End      time.Time `json:"dateTo"`
		LKSLogin string    `json:"staffNumber"`
	}

	payload.Start = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, &time.Location{})
	payload.End = payload.Start.AddDate(0, 1, -1)
	payload.LKSLogin = p.LksLogin

	b, err := json.Marshal(&payload)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("body marshal error: %e", err))
		return nil, err
	}

	req.SetBody(b)

	log.Debug().Msg("send request...")
	if err := fasthttp.DoTimeout(req, res, 20*time.Second); err != nil {
		log.Debug().Msg(fmt.Sprintf("sending request error: %e", err))
		return nil, err
	}

	if status := res.StatusCode(); status != 200 {
		if status == 302 {
			log.Debug().Msg(fmt.Sprintf("unauthorized error: %e", ErrExpiredCookie))
			return nil, ErrExpiredCookie
		}
		log.Debug().Msg(fmt.Sprintf("send request returs status %d", err))
		return nil, fmt.Errorf("send request return status %d", status)
	}

	var resBody monthPerspectiveResponse
	if err := json.Unmarshal(res.Body(), &resBody); err != nil {
		log.Debug().Msg(fmt.Sprintf("response body unmarshal error: %e", err))
		return nil, err
	}

	return resBody.Model.PerspectivePlan, nil

}

func (a *LksAPI) GetArchiveDuty(ctx context.Context, p AuthPayload, month int, year int) ([]ArchiveDuty, error) {
	log := a.getLogger("get_archive_duty")
	cookie, err := a.auth(ctx, p.AccordLogin, p.AccordPassword, p.LksLogin, p.LksPassword)
	if err != nil {
		return nil, err
	}

	log.Debug().Msg("prepare request")
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetRequestURI(archiveOrderUrl)
	req.Header.SetMethod(fasthttp.MethodPost)
	for name, value := range cookie {
		req.Header.SetCookie(name, value)
	}
	req.Header.SetContentType("application/json; charset=UTF-8")

	var payload struct {
		DateFrom   string `json:"dateFrom"`
		DateTo     string `json:"dateTo"`
		Count      int    `json:"itemsCount"`
		StartIndex int    `json:"startIndex"`
		LKSLogin   string `json:"staffNumber"`
	}

	dateFrom := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, &time.Location{})
	dateTo := dateFrom.AddDate(0, 1, -1)

	payload.DateFrom = dateFrom.Format("02.01.2006 15:04")
	payload.DateTo = dateTo.Format("02.01.2006 15:04")
	payload.LKSLogin = p.LksLogin
	payload.Count = 100
	payload.StartIndex = 0

	b, err := json.Marshal(&payload)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("body marshal error: %e", err))
		return nil, err
	}

	req.SetBody(b)

	log.Debug().Msg("send request...")
	if err := fasthttp.DoTimeout(req, res, 20*time.Second); err != nil {
		log.Debug().Msg(fmt.Sprintf("sending request error: %e", err))
		return nil, err
	}

	if status := res.StatusCode(); status != 200 {
		if status == 302 {
			log.Debug().Msg(fmt.Sprintf("unauthorized error: %e", ErrExpiredCookie))
			return nil, ErrExpiredCookie
		}
		log.Debug().Msg(fmt.Sprintf("send request returs status %d", err))
		return nil, fmt.Errorf("send request return status %d", status)
	}

	var resBody archiveDutyResponse
	if err := json.Unmarshal(res.Body(), &resBody); err != nil {
		log.Debug().Msg(fmt.Sprintf("response body unmarshal error: %e", err))
		return nil, err
	}
	return resBody.Model.AchievedDuties, nil
}

func (a *LksAPI) auth(ctx context.Context, accLogin string, accPass string, lksLogin string, lksPass string) (map[string]string, error) {
	log := a.getLogger("auth")
	log.Debug().Msg("try extract auth cookie from cache...")
	if c, ok := a.cache.Get(lksLogin); ok {
		if err := a.checkAuthCookie(c); err != nil {
			log.Debug().Msg(fmt.Sprintf("auth cookie invalid: %e", err))
			a.cache.Delete(lksLogin)
		} else {
			log.Debug().Msg("valid auth cookie has in cache")
			return c, nil
		}
		log.Debug().Msg("auth cookie has not in cache")
	}

	c, err := a.pool.Auth(ctx, accLogin, accPass, lksLogin, lksPass)
	if err != nil {
		return nil, err
	}

	if err := a.cache.Put(lksLogin, c); err != nil {
		log.Warn().Msg(fmt.Sprintf("put cookie in cache error: %e", err))
	}

	return c, nil
}

func (a *LksAPI) checkAuthCookie(c map[string]string) error {
	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetRequestURI(currentOrderUrl)
	req.Header.SetMethod(fasthttp.MethodPost)
	for name, value := range c {
		req.Header.SetCookie(name, value)
	}
	req.Header.SetContentType("application/json; charset=UTF-8")

	if err := fasthttp.DoTimeout(req, res, 20*time.Second); err != nil {
		return fmt.Errorf("send request error: %e", err)
	}

	switch res.StatusCode() {
	case 200:
		return nil
	default:
		return ErrExpiredCookie
	}
}

func (a *LksAPI) getLogger(method string) zerolog.Logger {
	return a.log.With().Str("method", method).Logger()
}
