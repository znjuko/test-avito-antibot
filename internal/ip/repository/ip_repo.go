package repository

import (
	"database/sql"
	"time"
)

type IpRepository struct {
	database *sql.DB
}

func NewIpRepository(db *sql.DB) IpRepository {
	return IpRepository{database: db}
}

func (Ip IpRepository) CreateFirstIpMask(ip string) error {
	_, err := Ip.database.Exec("INSERT INTO ip (ip_value,counter) VALUES ($1,$2)", ip, 1)

	return err
}

func (Ip IpRepository) GetMaskData(ip string) (int, string, error) {

	var ipCounter *int
	var timeout *string
	ipRow := Ip.database.QueryRow("SELECT counter , time FROM ip WHERE ip_value = $1", ip)

	err := ipRow.Scan(&ipCounter, &timeout)

	if err != nil {
		return 0, "", err
	}

	if timeout == nil {
		defaultTime := ""
		timeout = &defaultTime
	}

	return *ipCounter, *timeout, nil
}

func (Ip IpRepository) SetMaskDataToDefault(ip string) error {
	_, err := Ip.database.Exec("UPDATE ip SET counter = 1 , time = NULL where ip_value = $1", ip)

	return err
}

func (Ip IpRepository) UpdateMaskData(ip string) error {

	_, err := Ip.database.Exec("UPDATE ip SET counter = counter + 1  where ip_value = $1", ip)

	return err
}

func (Ip IpRepository) SetMaskCoolDown(ip string, limit int, coolDown time.Time) error {
	_, err := Ip.database.Exec("UPDATE ip SET counter = $2 , time = $3 where ip_value = $1", ip, limit, coolDown)

	return err
}

func (Ip IpRepository) ResetIpCoolDown(ip string) error {
	_, err := Ip.database.Exec("UPDATE ip SET counter = 0 , time = NULL where ip_value = $1", ip)

	return err
}
