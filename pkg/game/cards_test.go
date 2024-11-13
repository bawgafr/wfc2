package game

import (
	"reflect"
	"testing"
)

func Test_convertConnections(t *testing.T) {
	string := "GGRG"
	got := convertConnections(string)

	want := []Connector{Grass, Grass, Road, Grass}

	if len(want) != len(got) {
		t.Errorf("got %v, want %v", got, want)
	}

	for i := range want {
		if want[i] != got[i] {
			t.Errorf("got %v, want %v", got, want)
		}
	}

}

// test with an inital L shape, and rotate it 90, 180, 270 degrees
func Test_rotateConnection(t *testing.T) {
	initial := []Connector{Road, Road, Grass, Grass}

	t.Run("rotate 90", func(t *testing.T) {
		got := rotateConnections(initial, 90)
		want := []Connector{Grass, Road, Road, Grass}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("rotate 180", func(t *testing.T) {
		want := rotateConnections(initial, 180)
		got := []Connector{Grass, Grass, Road, Road}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("rotate 270", func(t *testing.T) {
		got := rotateConnections(initial, 270)
		want := []Connector{Road, Grass, Grass, Road}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

}
