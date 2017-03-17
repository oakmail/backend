package models

import (
	"time"
)

// PublicKey is the representation of a public key in the system
type PublicKey struct {
	ID           uint64    `db:"key_id" json:"key_id"`
	DateCreated  time.Time `db:"date_created" json:"date_created"`
	DateModified time.Time `db:"date_modified" json:"date_modified"`
	Owner        uint64    `db:"owner" json:"owner"`

	Algorithm        uint8  `db:"algorithm" json:"algorithm"`
	Length           uint16 `db:"length" json:"length"`
	Body             []byte `db:"body" json:"body"`
	KeyIDString      string `db:"key_id_string" json:"key_id_string"`
	KeyIDShortString string `db:"key_id_short_string" json:"key_id_short_string"`
	MasterKey        uint64 `db:"master_key" json:"master_key"`
}

// PublicKeyIdentity is part of the metadata of the key, a signed identity of the key
type PublicKeyIdentity struct {
	ID            uint64   `db:"identity" json:"identity" goqu:"skipinsert"`
	Key           uint64   `db:"key" json:"key"`
	Name          string   `db:"name" json:"name"`
	SelfSignature uint64   `db:"self_signature" json:"self_signature"`
	Signatures    []uint64 `db:"signatures" json:"signatures"`
}

// PublicKeySignature is a signature of an identity. Creates Web of Trust.
type PublicKeySignature struct {
	ID                   uint64    `db:"id" json:"id" goqu:"skipinsert"`
	Type                 uint8     `db:"type" json:"type"`
	Algorithm            uint8     `db:"algorithm" json:"algorithm"`
	Hash                 uint      `db:"hash" json:"hash"`
	CreationTime         time.Time `db:"creation_time" json:"creation_time"`
	SigLifetimeSecs      uint32    `db:"sig_lifetime_secs" json:"sig_lifetime_secs"`
	KeyLifetimeSecs      uint32    `db:"key_lifetime_secs" json:"key_lifetime_secs"`
	IssuerKeyID          uint64    `db:"issuer_key_id" json:"issuer_key_id"`
	IsPrimaryID          bool      `db:"is_primary_id" json:"is_primary_id"`
	RevocationReason     uint8     `db:"revocation_reason" json:"revocation_reason"`
	RevocationReasonText string    `db:"revocation_reason_text" json:"revocation_reason_text"`
}
