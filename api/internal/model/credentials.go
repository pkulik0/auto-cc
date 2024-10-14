package model

import (
	"strings"

	"github.com/pkulik0/autocc/api/internal/pb"
	"gorm.io/gorm"
)

func secureSecret(s string) string {
	lenVisible := len(s) / 3
	lenHidden := len(s) - lenVisible
	return s[:lenVisible] + strings.Repeat("*", lenHidden)
}

// CredentialsDeepL is a model for storing DeepL API credentials.
type CredentialsDeepL struct {
	gorm.Model
	Key   string
	Usage uint
}

// ToProto converts the model to a protobuf message.
func (c *CredentialsDeepL) ToProto() *pb.CredentialsDeepL {
	return &pb.CredentialsDeepL{
		Id:    uint64(c.ID),
		Key:   secureSecret(c.Key),
		Usage: uint64(c.Usage),
	}
}

// CredentialsGoogle is a model for storing Google API credentials.
type CredentialsGoogle struct {
	gorm.Model
	ClientID     string
	ClientSecret string
	Usage        uint
}

// ToProto converts the model to a protobuf message.
func (c *CredentialsGoogle) ToProto() *pb.CredentialsGoogle {
	return &pb.CredentialsGoogle{
		Id:           uint64(c.ID),
		ClientId:     c.ClientID,
		ClientSecret: secureSecret(c.ClientSecret),
		Usage:        uint64(c.Usage),
	}
}
