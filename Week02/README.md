## 作业

我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？



## 解答

在业务层面，sql.ErrNoRows和其他的error是不同的，因此在dao层应该提供区分它们的能力，不应该让上层与database/sql包耦合。所以，对于sql.ErrNoRows的error，通过「fmt.Errorf("%w")」的方式在dao层定义一个`sentinel error`  然后往上层抛，对于其他错误则直接wrap error然后抛给上层。



## 程序代码

1. 模拟查询用户信息的操作，表结构如下

| id   | name  | age  |
| ---- | ----- | ---- |
| 1    | Lili  | 18   |
| 2    | Tom   | 20   |
| 3    | Grace | 30   |

2. 启动http服务

```go
func main() {
	http.HandleFunc("/user", controller.GetUser)
	log.Println(http.ListenAndServe(":8080", nil))
}
```

3. controller层

```go
func GetUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Pares request error: ", err)
	}
	resp := &model.Response{}
	name := r.Form.Get("name")
	user, err := biz.FindUser(name)
	if err != nil {
		if errors.Is(err, dao.ErrNotFound) {
			log.Printf("ErrNoRows: %+v", err)
			w.WriteHeader(http.StatusNotFound)
			resp.Code = 404
			resp.Msg = "user not found"
			SendResp(w, resp)
			return
		} else {
			log.Printf("OtherError: %+v", err)
			w.WriteHeader(http.StatusInternalServerError)
			resp.Code = 500
			resp.Msg = "server not available"
			SendResp(w, resp)
			return
		}
	}
	resp.Code = 200
	resp.Msg = "success"
	resp.Data = user
	SendResp(w, resp)
}
```

4. biz层

```go
func FindUser(name string) (*model.User, error) {
	return service.FindUserByName(name)
}
```

5. service层

```go
func FindUserByName(name string) (*model.User, error) {
	user, err := dao.FindUserByName(name)
	if err != nil {
		return nil, err
	}
	user.Name += "是个好孩子！"
	return user, nil
}
```

6. dao层

```go
func FindUserByName(name string) (*model.User, error) {
	db, err := NewDB(dbName)
	if err != nil {
		return nil, err
	}
	userSQL := "SELECT * FROM userinfo WHERE `name` = ?"
	user := &model.User{}
	err = db.QueryRow(userSQL, name).Scan(&user.Id, &user.Name, &user.Age)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrapf(ErrNotFound, "sql [%s]", userSQL)
		}
		return nil, errors.Wrapf(err, "get user name[%s] error", name)
	}
	return user, nil
}
```



### [学习笔记](note.md)

