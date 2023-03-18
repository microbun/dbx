package dbx

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

var db *DB
var err error


type ExampleRecord struct {
	ID                 int64      `dbx:"column:id,primary_key,auto_increment" json:"id" `
	Tinyint            int8       `dbx:"column:tinyint" json:"tinyint" `
	Smallint           int16      `dbx:"column:smallint" json:"smallint" `
	Mediumint          int32      `dbx:"column:mediumint" json:"mediumint" `
	Integer            int32      `dbx:"column:integer" json:"integer" `
	Int                int32      `dbx:"column:int" json:"int" `
	Bigint             int64      `dbx:"column:bigint" json:"bigint" `
	Decimal            float64    `dbx:"column:decimal" json:"decimal" `
	Numeric            float64    `dbx:"column:numeric" json:"numeric" `
	Float              float32    `dbx:"column:float" json:"float" `
	Double             float64    `dbx:"column:double" json:"double" `
	Bit                []byte       `dbx:"column:bit" json:"bit" `
	Datetime           time.Time  `dbx:"column:datetime" json:"datetime" `
	Timestamp          time.Time  `dbx:"column:timestamp" json:"timestamp" `
	Char               string     `dbx:"column:char" json:"char" `
	Varchar            string     `dbx:"column:varchar" json:"varchar" `
	Enum               string     `dbx:"column:enum" json:"enum" `
	Bool               bool       `dbx:"column:bool" json:"bool" `
	NullableTinyint    *int8      `dbx:"column:nullable_tinyint" json:"nullable_tinyint,omitempty" `
	NullableSmallint   *int16     `dbx:"column:nullable_smallint" json:"nullable_smallint,omitempty" `
	NullableMediumint  *int32     `dbx:"column:nullable_mediumint" json:"nullable_mediumint,omitempty" `
	NullableInteger    *int32     `dbx:"column:nullable_integer" json:"nullable_integer,omitempty" `
	NullableInt        *int32     `dbx:"column:nullable_int" json:"nullable_int,omitempty" `
	NullableBigint     *int64     `dbx:"column:nullable_bigint" json:"nullable_bigint,omitempty" `
	NullableDecimal    *float64   `dbx:"column:nullable_decimal" json:"nullable_decimal,omitempty" `
	NullableNumeric    *float64   `dbx:"column:nullable_numeric" json:"nullable_numeric,omitempty" `
	NullableFloat      *float32   `dbx:"column:nullable_float" json:"nullable_float,omitempty" `
	NullableDouble     *float64   `dbx:"column:nullable_double" json:"nullable_double,omitempty" `
	NullableBit        *bool      `dbx:"column:nullable_bit" json:"nullable_bit,omitempty" `
	NullableDate       *time.Time `dbx:"column:nullable_date" json:"nullable_date,omitempty" `
	NullableDatetime   *time.Time `dbx:"column:nullable_datetime" json:"nullable_datetime,omitempty" `
	NullableTimestamp  *time.Time `dbx:"column:nullable_timestamp" json:"nullable_timestamp,omitempty" `
	NullableChar       *string    `dbx:"column:nullable_char" json:"nullable_char,omitempty" `
	NullableVarchar    *string    `dbx:"column:nullable_varchar" json:"nullable_varchar,omitempty" `
	NullableBinary     []byte     `dbx:"column:nullable_binary" json:"nullable_binary,omitempty" `
	NullableVarbinary  []byte     `dbx:"column:nullable_varbinary" json:"nullable_varbinary,omitempty" `
	NullableBlob       []byte     `dbx:"column:nullable_blob" json:"nullable_blob,omitempty" `
	NullableMediumblob []byte     `dbx:"column:nullable_mediumblob" json:"nullable_mediumblob,omitempty" `
	NullableLongtext   *string    `dbx:"column:nullable_longtext" json:"nullable_longtext,omitempty" `
	NullableMediumtext *string    `dbx:"column:nullable_mediumtext" json:"nullable_mediumtext,omitempty" `
	NullableText       *string    `dbx:"column:nullable_text" json:"nullable_text,omitempty" `
	NullableEnum       *string    `dbx:"column:nullable_enum" json:"nullable_enum,omitempty" `
	NullableBool       *bool      `dbx:"column:nullable_bool" json:"nullable_bool,omitempty" `
}


func (_ *ExampleRecord) TableName() string {
	return "example"
}

func (r *ExampleRecord) JSON() string {
	s, _ := json.Marshal(r)
	return string(s)
}

func (r *ExampleRecord) PettyJSON() string {
	s, _ := json.MarshalIndent(r, "", "\t")
	return string(s)
}

//	ID          int64           `dbx:"column:id,primary_key" `
//	CreatedAt   sql.NullTime    `dbx:"column:created_at,insert:current_timestamp"`
//	UpdatedAt   sql.NullTime    `dbx:"column:updated_at,update:current_timestamp,insert:current_timestamp"`
//	NullBool    sql.NullBool    `dbx:"column:null_bool" `
//	NullInt32   sql.NullInt32   `dbx:"column:null_int32"`
//	NullInt64   sql.NullInt64   `dbx:"column:null_int64"`
//	NullFloat64 sql.NullFloat64 `dbx:"column:null_float64"`
//	NullString  sql.NullString  `dbx:"column:null_string"`
//	NullTime    sql.NullTime    `dbx:"column:null_time"`
//	Int         int             `dbx:"column:int" `
//
//	Status    int            `dbx:"column:status" `
//	CreatedBy sql.NullString `dbx:"column:created_by"`
//	Avatar    sql.NullString `dbx:"column:avatar"`
//	NetworkID sql.NullInt64  `dbx:"column:network_id"`
//}

func TestMain(m *testing.M) {
	ExampleOpen()
	code := m.Run()
	fmt.Println("end testing")
	os.Exit(code)
}
