package decimal

import (
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestBSONRoundTrip(t *testing.T) {
	type doc struct {
		D Decimal `bson:"d"`
	}

	cases := []string{
		"0",
		"1",
		"-1",
		"123.456",
		"-987654321.123456789",
		"12345678901234567890123456789012345.6789",
		"1e100",
		"-1e-100",
	}

	for _, s := range cases {
		t.Run(s, func(t *testing.T) {
			in, err := NewFromString(s)
			if err != nil {
				t.Fatalf("NewFromString(%q): %v", s, err)
			}
			data, err := bson.Marshal(doc{D: in})
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			var out doc
			if err := bson.Unmarshal(data, &out); err != nil {
				t.Fatalf("Unmarshal: %v", err)
			}
			if !in.Equal(out.D) {
				t.Fatalf("round-trip mismatch: in=%s out=%s", in, out.D)
			}
		})
	}
}

func TestBSONUnmarshalEmptyString(t *testing.T) {
	type doc struct {
		D Decimal `bson:"d"`
	}
	data, err := bson.Marshal(struct {
		D string `bson:"d"`
	}{D: ""})
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var out doc
	if err := bson.Unmarshal(data, &out); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !out.D.Equal(Decimal{}) {
		t.Fatalf("expected zero Decimal, got %s", out.D)
	}
}
