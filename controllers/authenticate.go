package controllers

import (
	"crypto"
	"fmt"
	"log"
	"time"

	"github.com/bytemare/ksf"
	"github.com/bytemare/opaque"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var OpaqueConfig = &opaque.Configuration{
	OPRF: opaque.OPRFP256,           // Oblivious Pseudorandom Function
	KDF:  crypto.SHA512,             // Key Derivation Function
	MAC:  crypto.SHA512,             // Message Authentication Code
	Hash: crypto.SHA512,             // Hash function
	KSF:  ksf.Argon2id,              // Key Stretching Function (e.g., Argon2id)
	AKE:  opaque.Ristretto255Sha512, // Authenticated Key Exchange
}

func InitConfig() {
	if OpaqueConfig == nil {
		log.Fatal("Failed to initialize OpaqueConfig")
	}
	log.Println("OpaqueConfig initialized successfully")
}

// RegisterUser registers a new user, creating an OPAQUE envelope for secure authentication
func RegisterUser(username, email string, profilePicture []byte, password string, db *sqlx.DB) error {
	// Step 1: OPAQUE Client-Side Initiation
	clientSession, msg1, err := opaque.PwRegInit(username, password, 4096)
	if err != nil {
		return fmt.Errorf("error initializing OPAQUE password registration: %v", err)
	}

	// Step 2: Server-Side Response to Client Msg1 using OpaqueConfig
	_, msg2, err := opaque.PwReg1(OpaqueConfig, msg1)
	if err != nil {
		return fmt.Errorf("error processing OPAQUE password registration: %v", err)
	}

	// Step 3: Client Finalization Using Msg2
	msg3, err := opaque.PwReg2(clientSession, msg2)
	if err != nil {
		return fmt.Errorf("OPAQUE client finalization failed: %v", err)
	}

	// Step 4: Storing the User Information in PostgreSQL
	// Use a UUID for the user ID and store other required fields in the users table
	query := `
		INSERT INTO users (id, username, email, profile_picture, role, opaque_envelope, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	userID := uuid.New()
	_, err = db.Exec(
		query,
		userID,           // id
		username,         // username
		email,            // email
		profilePicture,   // profile_picture (if null, pass `nil`)
		"user",           // role
		msg3.Serialize(), // opaque_envelope
		time.Now().UTC(), // created_at
	)
	if err != nil {
		return fmt.Errorf("error inserting user into database: %v", err)
	}

	log.Printf("User %s successfully registered with ID %s", username, userID.String())
	return nil
}

// TODO //  - Implement LOGIN FUNCTIONALITY
