package auth

import "testing"

func TestPasswordHashRoundTrip(t *testing.T) {
	hash, err := HashPassword("secret")
	if err != nil {
		t.Fatal(err)
	}
	if !CheckPassword(hash, "secret") {
		t.Fatal("correct password rejected")
	}
	if CheckPassword(hash, "wrong") {
		t.Fatal("wrong password accepted")
	}
}

func TestTokenRoundTrip(t *testing.T) {
	tok, err := GenerateToken("user-1", "a@b.com", "s3cret", 0)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := ParseToken(tok, "s3cret")
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != "user-1" || claims.Email != "a@b.com" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
	if _, err := ParseToken(tok, "wrong"); err == nil {
		t.Fatal("expected error for wrong secret")
	}
}
