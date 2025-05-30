package h264_test

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/tachode/bitstream-go/h264"
)

var parser = h264.NewParser()

func TestParse_ValidSPS(t *testing.T) {
	spsNalUnit := spsBytes[:]

	nal, err := parser.Parse(spsNalUnit)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if nal.NalUnitType != h264.NalUnitTypeSPS {
		t.Errorf("Expected NalUnitTypeSPS, got %v", nal.NalUnitType)
	}

	if _, ok := nal.Payload.(*h264.SeqParameterSet); !ok {
		t.Errorf("Expected payload to be of type SeqParameterSet")
	}

	nalJSON, err := json.MarshalIndent(nal, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal nal to JSON: %v", err)
	}

	t.Logf("NAL JSON:\n%s", nalJSON)
}

func TestParse_ValidPPS(t *testing.T) {
	ppsNalUnit := ppsBytes[:]

	nal, err := parser.Parse(ppsNalUnit)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if nal.NalUnitType != h264.NalUnitTypePPS {
		t.Errorf("Expected NalUnitTypePPS, got %v", nal.NalUnitType)
	}

	if _, ok := nal.Payload.(*h264.PicParameterSet); !ok {
		t.Errorf("Expected payload to be of type SeqParameterSet")
	}

	nalJSON, err := json.MarshalIndent(nal, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal nal to JSON: %v", err)
	}

	t.Logf("NAL JSON:\n%s", nalJSON)
}

func TestParse_ValidSEI(t *testing.T) {
	for i, seiNalUnit := range seiBytes {
		t.Logf("SEI test %d\n%s", i, hex.Dump(seiNalUnit))
		nal, err := parser.Parse(seiNalUnit)
		if err != nil {
			t.Fatalf("Parse failed: %v", err)
		}

		if nal.NalUnitType != h264.NalUnitTypeSEI {
			t.Errorf("Expected NalUnitTypeSEI, got %v", nal.NalUnitType)
		}

		if _, ok := nal.Payload.(*h264.Sei); !ok {
			t.Errorf("Expected payload to be of type Sei")
		}

		nalJSON, err := json.MarshalIndent(nal, "", "  ")
		if err != nil {
			t.Fatalf("Failed to marshal nal to JSON: %v", err)
		}

		t.Logf("NAL JSON:\n%s", nalJSON)
	}
}

func TestParse_EmptyBuffer(t *testing.T) {
	// Empty buffer
	emptyBuffer := []byte{}

	_, err := parser.Parse(emptyBuffer)
	if err == nil {
		t.Fatalf("Expected error for empty buffer, got nil")
	}
}

func TestParse_ValueJson(t *testing.T) {
	t.Logf("Extracted Values:\n%s", parser.ValueJson())
}
