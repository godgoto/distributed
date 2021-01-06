package db_token

import (
)

type MToken struct {
 Id int `gorm:auto_increment`
 Token string `gorm:varchar(64)`
 }