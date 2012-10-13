package llsd

import (
	"testing"
)

func TestSimpleArray(t *testing.T) {
	xml := []byte(`
<llsd>
    <array>
        <integer>1</integer>
        <integer>2</integer>
        <integer>3</integer>
    </array>
</llsd>
    `)
	var v Array
	err := UnmarshalXML(xml, &v)
	if err != nil {
		t.Fatalf("Error unmarshalling XML: %s", err)
	}
	t.Logf("Array values: %v", v)
	if len(v) != 3 {
		t.Errorf("Expected XML length: 3; actual length: %d", len(v))
	}
	if v[0] != 1 || v[1] != 2 || v[2] != 3 {
		t.Errorf("Array has wrong values: %v", v)
	}
}

func TestSimpleMap(t *testing.T) {
	xml := []byte(`
<llsd>
    <map>
        <key>One</key>
        <integer>1</integer>
        <key>Two</key>
        <integer>2</integer>
        <key>Three</key>
        <integer>3</integer>
    </map>
</llsd>
    `)
	var v Map
	err := UnmarshalXML(xml, &v)
	if err != nil {
		t.Fatalf("Error unmarshalling XML: %s", err)
	}
	t.Logf("Map values: %v", v)
	if len(v) != 3 {
		t.Errorf("Expected map of size 3, got %d", len(v))
	}
	if v["One"] != 1 || v["Two"] != 2 || v["Three"] != 3 {
		t.Errorf("Map has wrong values: %v", v)
	}
}

func TestBoolean(t *testing.T) {
	xml := []byte(`
<llsd>
    <array>
        <boolean>true</boolean>
        <boolean>1</boolean>
        <boolean>false</boolean>
        <boolean>0</boolean>
        <boolean />
    </array>
</llsd>
    `)
	var v Array
	err := UnmarshalXML(xml, &v)
	t.Logf("Booleans: %v", v)
	if err != nil {
		t.Fatalf("Error unmarshalling XML: %s", err)
	}
	if !v[0].(bool) {
		t.Error("'true' is not true")
	}
	if !v[1].(bool) {
		t.Error("1 is not true")
	}
	if v[2].(bool) {
		t.Error("'false is not false")
	}
	if v[3].(bool) {
		t.Error("'0' is not false")
	}
	if v[4].(bool) {
		t.Error("<boolean /> is not false")
	}
}

func TestDummyAsset(t *testing.T) {
	xml := []byte(`
<?xml version="1.0" encoding="UTF-8"?>
<llsd>
    <array>
        <map>
            <key>creation-date</key>
            <date>2007-03-15T18:30:18Z</date>
            <key>creator-id</key>
            <uuid>3c115e51-04f4-523c-9fa6-98aff1034730</uuid>
        </map>
        <string>0123456789</string>
        <string>Where's the beef?</string>
        <string>Over here.</string>
        <string>default
{
    state_entry()
    {
        llSay(0, "Hello, Avatar!");
    }

    touch_start(integer total_number)
    {
        llSay(0, "Touched.");
    }
}</string>
        <binary encoding='base64'>AABAAAAAAAAAAAIAAAA//wAAP/8AAADgAAAA5wAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABkAAAAZAAAAAAAAAAAAAAAZAAAAAAAAAABAAAAAAAAAAAAAAAAAAAABQAAAAEAAAAQAAAAAAAAAAUAAAAFAAAAABAAAAAAAAAAPgAAAAQAAAAFAGNbXgAAAABgSGVsbG8sIEF2YXRhciEAZgAAAABcXgAAAAhwEQjRABeVAAAABQBjW14AAAAAYFRvdWNoZWQuAGYAAAAAXF4AAAAIcBEI0QAXAZUAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA</binary>
    </array>
</llsd>
    `)
	var v Array
	err := UnmarshalXML(xml, &v)
	if err != nil {
		t.Fatalf("Error unmarshalling XML: %s", err)
	}

	if len(v) != 6 {
		t.Error("Array has incorrect length.")
	}
	_, ok := v[0].(Map)
	if !ok {
		t.Error("First element is not a map.")
	}
}

func TestUndefRootArray(t *testing.T) {
	xml := []byte(`<llsd><undef /></llsd>`)
	var v Array
	err := UnmarshalXML(xml, &v)
	if err != nil {
		t.Fatalf("Failed to unmarshal undef into an array: %s", err)
	}
	if len(v) != 0 {
		t.Error("Undefined array is not of length zero.")
	}
}

func TestUndefRootMap(t *testing.T) {
	xml := []byte(`<llsd><undef /></llsd>`)
	var v Map
	err := UnmarshalXML(xml, &v)
	if err != nil {
		t.Fatalf("Failed to unmarshal undef into a map: %s", err)
	}
	if len(v) != 0 {
		t.Error("Undefined map is not of length zero.")
	}
}

func TestUndefInterface(t *testing.T) {
	xml := []byte(`<llsd><undef /></llsd>`)
	var v interface{}
	err := UnmarshalXML(xml, &v)
	if err != nil {
		t.Fatalf("Failed to unmarshal undef into an interface: %s", err)
	}
}
