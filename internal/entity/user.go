package entity

import "time"

type GroupID string
type UserState string

const (
	Group1_1 GroupID = "1.1"
	Group1_2 GroupID = "1.2"
	Group2_1 GroupID = "2.1"
	Group2_2 GroupID = "2.2"
	Group3_1 GroupID = "3.1"
	Group3_2 GroupID = "3.2"
	Group4_1 GroupID = "4.1"
	Group4_2 GroupID = "4.2"
	Group5_1 GroupID = "5.1"
	Group5_2 GroupID = "5.2"
	Group6_1 GroupID = "6.1"
	Group6_2 GroupID = "6.2"
)

const (
	Idle UserState = "Idle"
)

type User struct {
	ID         string
	TelegramID int64
	Username   string
	FirstName  string
	Groups     []GroupID
	IsActive   bool
	State      UserState
	CreatedAt  time.Time
}
