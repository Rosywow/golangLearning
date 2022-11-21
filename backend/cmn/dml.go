package cmn

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

func jsonDataTypeToGo(v interface{}) {
	//just silent return
	if v == nil {
		return
	}

	//process array
	switch reflect.TypeOf(v).Kind() {
	case reflect.Slice, reflect.Array:
		array, _ := v.([]interface{})
		for i, e := range array {
			switch reflect.TypeOf(e).Kind() {
			case reflect.Slice, reflect.Array:
				jsonDataTypeToGo(e)
			case reflect.Float32, reflect.Float64:
				data, _ := e.(float64)
				_, frac := math.Modf(data)
				if frac == 0 {
					array[i] = int64(data)
				}
			case reflect.Map:
				jsonDataTypeToGo(e)
			}
		}
		return
	}

	m, ok := v.(map[string]interface{})
	if !ok {
		return
	}

	for k, v := range m {
		switch reflect.TypeOf(v).Kind() {
		case reflect.Float32, reflect.Float64:
			data, _ := v.(float64)
			_, frac := math.Modf(data)
			if frac == 0 {
				m[k] = int64(data)
			}

		case reflect.Slice, reflect.Array:
			array, _ := v.([]interface{})
			for i, e := range array {
				switch reflect.TypeOf(e).Kind() {
				case reflect.Float32, reflect.Float64:
					data, _ := e.(float64)
					_, frac := math.Modf(data)
					if frac == 0 {
						array[i] = int64(data)
					}
				case reflect.Map:
					jsonDataTypeToGo(e)
				}
			}

		case reflect.Map:
			jsonDataTypeToGo(v)
		}
	}
}

//DML Data Management Language
func DML(f *Filter, req *ReqProto) (err error) {
	if sqlxDB == nil {
		err = fmt.Errorf("please connect to dbms")
		z.Error(err.Error())
		return
	}
	if f == nil {
		err = fmt.Errorf("call executeDML with nil param f")
		z.Error(err.Error())
		return
	}
	if f.TableMap == nil {
		err = fmt.Errorf("call executeDML with nil param f.tableMap")
		z.Error(err.Error())
		return
	}
	if req.Action == "" {
		err = fmt.Errorf("call executeDML with empty param action")
		z.Error(err.Error())
		return
	}

	action := strings.ToUpper(req.Action)

	if reflect.TypeOf(reflect.ValueOf(f.TableMap).Interface()).Kind() != reflect.Ptr {
		err = fmt.Errorf("f.tableMap should be pointer, please using &struct to set it")
		z.Error(err.Error())
		return
	}

	objT := reflect.TypeOf(reflect.ValueOf(f.TableMap).Elem().Interface())

	fnc := reflect.ValueOf(f.TableMap).MethodByName("GetTableName")
	if !fnc.IsValid() {
		err = fmt.Errorf("missing GetTableName on struct " + objT.Name())
		z.Error(err.Error())
		return
	}
	r := fnc.Call([]reflect.Value{})
	if len(r) == 0 {
		err = fmt.Errorf("GetTableName on struct " + objT.Name() + " return empty")
		z.Error(err.Error())
		return
	}

	tblName, ok := r[0].Interface().(string)
	if !ok {
		err = fmt.Errorf("GetTableName on struct " + objT.Name() + " return non string value")
		z.Error(err.Error())
		return
	}

	if req.Filter != nil {
		jsonDataTypeToGo(req.Filter)
	}
	if req.AuthFilter != nil {
		jsonDataTypeToGo(req.AuthFilter)
	}

	//Query data prepare
	switch action {
	case "UPDATE", "INSERT":
		if len(req.Data) > 0 {
			err = json.Unmarshal(req.Data, f.TableMap)
			if err != nil {
				z.Error(err.Error())
				return
			}
		}

		// from pointer of struct to struct it self.
		objV := reflect.ValueOf(reflect.ValueOf(f.TableMap).Elem().Interface())

		// generate columns set and values set
		fieldNum := objV.NumField()

		// *** ATTENTION: we only support null.types
	nextField: // for continue using this label
		for i := 0; i < fieldNum; i++ {
			fieldType := objT.Field(i)
			fieldValue := objV.Field(i)
			if !fieldValue.CanInterface() { //non-exported field
				continue
			}

			var ok bool
			var dbTag string
			kind := fieldValue.Kind()

			dbTag, ok = fieldType.Tag.Lookup("db")
			if !ok { // not a db column field
				continue
			}
			// --------------
			// value check
			for {
				switch {
				case kind == reflect.Struct && fieldValue.FieldByName("Valid").IsValid():
					var valid bool
					valid, ok = fieldValue.FieldByName("Valid").Interface().(bool)
					if !ok {
						err = fmt.Errorf("valid field's data type should be bool, or you not using null.types")
						z.Error(err.Error())
						return
					}
					if !valid { // empty/nil field
						continue nextField
					}
				case (kind == reflect.Slice || kind == reflect.Array ||
					kind == reflect.Map) && fieldValue.Len() == 0:
					continue nextField
				}

				switch val := fieldValue.Interface().(type) {
				case bool:
					if !val {
						continue nextField
					}
				case int, int8, int16, int32, int64, float32, float64, uint, uint16, uint32, uint64:
					if val == 0 {
						continue nextField
					}
				case string:
					if val == "" {
						continue nextField
					}
				default:
					// it's struct/slice/array
					break
				}

				break
			}
			// ------------------

			columnName := strings.Trim(strings.Split(dbTag, ",")[0], " ")
			if action == "UPDATE" && strings.ToLower(columnName) == "id" {
				z.Info("id update disabled")
				continue
			}
			f.Columns = append(f.Columns, columnName)
			f.Values = append(f.Values, fieldValue.Interface())
		}
		jsonDataTypeToGo(f.Values)

	case "SELECT":
		if len(req.Sets) != 0 {
			break
		}
		fnc := reflect.ValueOf(f.TableMap).MethodByName("Fields")
		if !fnc.IsValid() {
			err = fmt.Errorf("missing Fields on struct " + objT.Name())
			z.Error(err.Error())
			return
		}
		r := fnc.Call([]reflect.Value{})
		if len(r) == 0 {
			err = fmt.Errorf("Fields on struct " + objT.Name() + " return empty")
			z.Error(err.Error())
			return
		}

		fields, ok := r[0].Interface().([]string)
		if !ok {
			err = fmt.Errorf("Fields on struct " + objT.Name() + " return non []string value")
			z.Error(err.Error())
			return
		}
		req.Sets = fields
	}

	// Begin DML process
	switch action {
	case "INSERT":
		if len(f.Columns) == 0 {
			err = fmt.Errorf("empty column list")
			z.Error(err.Error())
			return
		}
		if len(f.Columns) != len(f.Values) {
			err = fmt.Errorf("number of columns not equal to number of values")
			z.Error(err.Error())
			return
		}
		var columns string
		var values string
		for i := 0; i < len(f.Columns); i++ {
			columns = columns + f.Columns[i] + ","
			values = values + fmt.Sprintf("$%d,", i+1)
		}
		columns = columns[:len(columns)-1]
		values = values[:len(values)-1]

		s := fmt.Sprintf(`INSERT INTO %s(%s) VALUES(%s) RETURNING ID`, tblName, columns, values)
		var stmt *sql.Stmt
		stmt, err = sqlxDB.Prepare(s)
		if err != nil {
			z.Error(err.Error())
			return
		}
		defer stmt.Close()

		z.Info(s)
		z.Info(fmt.Sprintf("%v", f.Values))

		r := stmt.QueryRow(f.Values...)
		var id int64
		err = r.Scan(&id)
		if err != nil {
			z.Error(err.Error())
			return
		}
		if id > 0 {
			z.Info(fmt.Sprintf("insert successfully with id = %d", id))
			f.QryResult = id
		}

	case "UPDATE":
		if len(f.Columns) == 0 {
			f.QryResult = 1
			z.Info("len(f.column)==0, update needless.")
			return
		}

		var expr string
		expr, err = f.CreateFilter(req)
		if err != nil {
			return
		}

		var sets string
		for i := 0; i < len(f.Columns); i++ {
			sets = sets + fmt.Sprintf("%s=$%d,", f.Columns[i], i+1)
		}
		sets = sets[:len(sets)-1]
		if f.AuthExpr != "" {
			expr = fmt.Sprintf("(%s) and (%s)", expr, f.AuthExpr)
			for i := 0; i < len(f.AuthWhereValues); i++ {
				f.Values = append(f.Values, f.AuthWhereValues[i])
			}
		}

		s := fmt.Sprintf("UPDATE %s SET %s WHERE %s", tblName, sets, expr)
		var stmt *sql.Stmt
		stmt, err = sqlxDB.Prepare(s)
		if err != nil {
			z.Error(err.Error())
			return
		}
		defer stmt.Close()

		var result sql.Result
		z.Info(s)
		z.Info(fmt.Sprintf("%v", f.Values))
		result, err = stmt.Exec(f.Values...)
		if err != nil {
			z.Error(err.Error())
			return
		}
		var d int64
		if d, err = result.RowsAffected(); err != nil {
			z.Error(err.Error())
			return
		}

		f.QryResult = d
		if d > 0 {
			z.Info("update success")
		}

	case "DELETE":
		var expr string
		expr, err = f.CreateFilter(req)
		if err != nil {
			z.Error(err.Error())
			return
		}

		if f.AuthExpr != "" {
			expr = fmt.Sprintf("(%s) and (%s)", expr, f.AuthExpr)
			for i := 0; i < len(f.AuthWhereValues); i++ {
				f.Values = append(f.Values, f.AuthWhereValues[i])
			}
		}

		s := fmt.Sprintf("DELETE FROM %s WHERE %s", tblName, expr)
		z.Info(s)
		z.Info(fmt.Sprintf("%v", f.Values))
		var stmt *sql.Stmt
		stmt, err = sqlxDB.Prepare(s)
		if err != nil {
			z.Error(err.Error())
			return
		}
		defer stmt.Close()

		var result sql.Result
		result, err = stmt.Exec(f.Values...)
		if err != nil {
			z.Error(err.Error())
			return
		}
		var d int64
		if d, err = result.RowsAffected(); err != nil {
			z.Error(err.Error())
			return
		}

		f.QryResult = d
		if d > 0 {
			z.Info("delete success")
			return
		}

	case "SELECT":
		if len(req.Sets) == 0 {
			err = fmt.Errorf("empty select sets  on " + objT.Name())
			z.Error(err.Error())
			return
		}
		var pageExpr string
		if req.PageSize != 0 {
			pageExpr = fmt.Sprintf("LIMIT %d OFFSET %d", req.PageSize, req.Page*req.PageSize)
			z.Info(pageExpr)
		}

		var pkList []string
		pkList, err = f.getPrimaryKeys(true)
		if err != nil {
			z.Error(err.Error())
			return
		}

		if len(pkList) == 0 && strings.ToLower(tblName[:2]) != "v_" {
			err = fmt.Errorf("missing primary key on " + tblName)
			z.Error(err.Error())
			return
		}

		orderBy := strings.Join(pkList, ",")
		var orderByList []string
		for i := 0; i < len(req.OrderBy); i++ {
			m := req.OrderBy[i]
			if len(m) == 0 {
				continue
			}
			for key, value := range m {
				columnName, found := f.mapKey(key)
				if !found {
					err = fmt.Errorf("unknown " + key + " on " + tblName)
					z.Error(err.Error())
					return
				}
				value = strings.ToUpper(value)
				if value != "ASC" && value != "DESC" {
					err = fmt.Errorf("unknown " + value + " order type with " + key + " on " + tblName)
					z.Error(err.Error())
					return
				}
				orderByList = append(orderByList, columnName+" "+value)
			}
		}
		if len(orderByList) > 0 {
			orderBy = strings.Join(orderByList, ",")
		}

		var expr string
		expr, err = f.CreateFilter(req)
		if err != nil {
			return
		}
		if expr == "" {
			expr = "1=1"
		}

		var columns []string
		for _, k := range req.Sets {
			fieldName, jsonOPr := getColumnName(k)
			columnName, found := f.mapKey(fieldName)
			if !found {
				err = fmt.Errorf("unknown " + fieldName + " on " + tblName)
				z.Error(err.Error())
				return
			}
			if jsonOPr {
				if strings.Contains(k, "->>") {
					err = fmt.Errorf("%s cause the result can't be marshal to json, please use '->' replace '->>'", k)
					z.Error(err.Error())
					return
				}
				columnName = strings.ReplaceAll(k, fieldName, columnName) + " as " + columnName
			}
			columns = append(columns, columnName)
		}
		sets := strings.Join(columns, ",")
		// ---------
		// get row count
		if f.AuthExpr != "" {
			expr = fmt.Sprintf("(%s) and (%s)", expr, f.AuthExpr)
			for i := 0; i < len(f.AuthWhereValues); i++ {
				f.Values = append(f.Values, f.AuthWhereValues[i])
				//f.whereValues = append(f.whereValues, f.authWhereValues[i])
			}
		}
		s := fmt.Sprintf("SELECT count(*) as row_count FROM %s WHERE %s", tblName, expr)
		z.Info(s)
		var stmt *sqlx.Stmt

		stmt, err = sqlxDB.Preparex(s)
		if err != nil {
			z.Error(err.Error())
			return
		}
		defer stmt.Close()

		row := stmt.QueryRowx(f.Values...)
		err = row.Scan(&f.RowCount)
		if err != nil {
			z.Error(err.Error())
			return
		}

		// ---------
		s = fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY %s %s", sets, tblName, expr, orderBy, pageExpr)
		z.Info(s)
		z.Info(fmt.Sprintf("%v", f.Values))
		stmt, err = sqlxDB.Preparex(s)
		if err != nil {
			z.Error(err.Error())
			return
		}
		defer stmt.Close()
		var rows *sqlx.Rows
		rows, err = stmt.Queryx(f.Values...)
		if err != nil {
			z.Error(err.Error())
			return
		}
		defer rows.Close()

		var buf []byte
		var qryResults []string
		for rows.Next() {
			pValue := reflect.New(objT)
			if !pValue.CanInterface() {
				err = fmt.Errorf("pValue can't interface() while it should")
				z.Error(err.Error())
				return
			}
			p := pValue.Interface()
			err = rows.StructScan(p)
			if err != nil {
				z.Error(err.Error())
				return
			}

			buf, err = MarshalJSON(p)
			if err != nil {
				z.Error(err.Error())
				return
			}
			var qryValue map[string]interface{}
			err = json.Unmarshal(buf, &qryValue)
			if err != nil {
				z.Error(err.Error())
				return
			}
			f.Result = append(f.Result, p)
			qryResults = append(qryResults, string(buf))
		}

		if len(qryResults) > 0 {
			f.QryResult = "[" + strings.Join(qryResults, ",") + "]"
			return
		}
		f.QryResult = "[]"

	default:
		err = fmt.Errorf("unsupported action " + action + " on " + objT.Name())
		z.Error(err.Error())
		return
	}
	return
}
