package code

import (
	"errors"
	"log"
	"net"
	"strings"
	"time"

	"github.com/bbconfhq/mycodes/models"
	"github.com/bbconfhq/mycodes/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type Response struct {
	Data  interface{}
	Error int
}

var (
	Validator = validator.New()
)

func GetList(c *fiber.Ctx) error {
	status, errs := fiber.StatusOK, 0

	result, err := repository.Code.GetRecent()
	if err != nil {
		status, errs = fiber.StatusInternalServerError, 1
		result = nil
		log.Printf("GetList error: %v", err)
	}

	return c.Status(status).JSON(Response{
		Data:  result,
		Error: errs,
	})
}

func GetOne(c *fiber.Ctx) error {
	status, errs := fiber.StatusOK, 0

	result, err := repository.Code.Get(c.Params("uid"))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		status, errs = fiber.StatusNotFound, 1
		result = nil
	} else if err != nil {
		status, errs = fiber.StatusInternalServerError, 1
		result = nil
		log.Printf("GetOne error: %v", err)
	}

	return c.Status(status).JSON(Response{
		Data:  result,
		Error: errs,
	})
}

func Post(c *fiber.Ctx) error {
	status, errs := fiber.StatusOK, 0
	result := &models.Code{}

	if err := c.BodyParser(result); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Data:  nil,
			Error: 1,
		})
	}

	if err := Validator.Struct(result); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(Response{
			Data:  nil,
			Error: 2,
		})
	}

	nid, err := gonanoid.New(10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(Response{
			Data:  nil,
			Error: 3,
		})
	}

	// The leftmost X-Forwarded-For entries are supplied by the client and can
	// be spoofed, so use the address appended by the nearest proxy
	rawIp := c.IP()
	if ips := c.IPs(); len(ips) > 0 {
		rawIp = strings.TrimSpace(ips[len(ips)-1])
	}
	// Keep the /16 of an IPv4 address or the /48 of an IPv6 address
	maskedIp := "0.0.0.0"
	if ip := net.ParseIP(rawIp); ip != nil {
		if v4 := ip.To4(); v4 != nil {
			maskedIp = v4.Mask(net.CIDRMask(16, 32)).String()
		} else {
			maskedIp = ip.Mask(net.CIDRMask(48, 128)).String()
		}
	}

	result = &models.Code{
		ID:        nid,
		Ip:        maskedIp,
		Name:      result.Name,
		Title:     result.Title,
		Content:   result.Content,
		Language:  result.Language,
		ExpiredAt: time.Now().AddDate(0, 0, 7),
	}

	err = repository.Code.Create(result)
	if err != nil {
		status, errs = fiber.StatusInternalServerError, 4
		result = nil
		log.Printf("Post error: %v", err)
	}

	return c.Status(status).JSON(Response{
		Data:  result,
		Error: errs,
	})
}
