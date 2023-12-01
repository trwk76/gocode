package postgres

import "github.com/trwk76/gocode"

type (
	ColumnType interface {
		gocode.Writable
		colType()
	}

	BooleanType  struct{}
	SmallIntType struct{}
	IntegerType  struct{}
	BigIntType   struct{}
	RealType     struct{}
	DoubleType   struct{}
	DateType     struct{}
	TextType     struct{}
	ByteaType    struct{}

	NumericType struct {
		Precision byte
		Scale     byte
	}

	TimeType struct {
		Precision uint8
		Timezone  bool
	}

	TimestampType struct {
		Precision uint8
		Timezone  bool
	}

	IntervalType struct {
		Precision uint8
	}

	CharType struct {
		Length uint32
	}

	VarCharType struct {
		Length uint32
	}

	VarType interface {
		ColumnType
		varType()
	}

	RecordType struct{}
)

func (BooleanType) colType()   {}
func (SmallIntType) colType()  {}
func (IntegerType) colType()   {}
func (BigIntType) colType()    {}
func (RealType) colType()      {}
func (DoubleType) colType()    {}
func (NumericType) colType()   {}
func (DateType) colType()      {}
func (TimeType) colType()      {}
func (TimestampType) colType() {}
func (IntervalType) colType()  {}
func (CharType) colType()      {}
func (VarCharType) colType()   {}
func (TextType) colType()      {}
func (ByteaType) colType()     {}
func (RecordType) colType()    {}
func (RecordType) varType()    {}

var (
	_ ColumnType = BooleanType{}
	_ ColumnType = SmallIntType{}
	_ ColumnType = IntegerType{}
	_ ColumnType = BigIntType{}
	_ ColumnType = RealType{}
	_ ColumnType = DoubleType{}
	_ ColumnType = NumericType{}
	_ ColumnType = DateType{}
	_ ColumnType = TimeType{}
	_ ColumnType = TimestampType{}
	_ ColumnType = IntervalType{}
	_ ColumnType = CharType{}
	_ ColumnType = VarCharType{}
	_ ColumnType = TextType{}
	_ ColumnType = ByteaType{}
	_ VarType    = RecordType{}
)
