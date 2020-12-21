
package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)
const (
	userName = "root"
	password = "123456"
	ip = "localhost"
	port = "3060"
	dbName = "users"
)
var DB *sql.DB
// 连接数据库
func initDB() {
	// 构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{userName, ":", password, "@tcp(",ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	// 打开数据库 （"驱动名"，连接）
	DB,_ = sql.Open("mysql",path)
	// 设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	// 设置数据库最大闲置数
	DB.SetMaxIdleConns(10)
	// 验证连接
	if err := DB.Ping(); err != nil {
		fmt.Println("连接数据库失败")
		return
	}
	fmt.Println("连接数据库成功")
}
type User struct {
	id int
	username string
	password string
}
func userLogin(c *gin.Context) {
	userName := c.Request.URL.Query().Get("username") //获取登录用户名
	passWord := c.Request.URL.Query().Get("password") //获取登录密码
	//查询列表
	rows,err := DB.Query("SELECT * FROM voro_user") //执行sql语句查询表中数据
	if err != nil {
		fmt.Println("查询失败")
	}
	var s User
	for rows.Next() {
		err = rows.Scan(&s.id, &s.username, &s.password)
		if err != nil {
			fmt.Println(err)
		}
	}
	if userName != s.username { // 判断获取到的登录用户名是否在数据库中存在
		// 无此用户(用户名不存在)
		c.JSON(200,gin.H{
			"success":false,
			"code":400,
			"msg":"无此用户",
		})
	} else {
		// 获取当前用户名密码
		// 获取登录用户名的密码 查询是否匹配
		us,_ := DB.Query("SELECT password FROM voro_user where username='" + userName + "'")
		for us.Next(){
			var u passWd
			err = us.Scan(&u.password)
			if err != nil {
				fmt.Println(err)
			}
			// 密码是否匹配
			if passWord != u.password{ //密码不一致
				c.JSON(200,gin.H{
					"success":false,
					"code":400,
					"msg":"密码错误",
				})
			} else {
				c.JSON(200,gin.H{ //用户名存在且密码匹配
					"success":true,
					"code":200,
					"msg":"登录成功",
				})
			}
		}
	}
	rows.Close()
}
func userRegister(c *gin.Context){
	userName := c.Request.URL.Query().Get("username")
	passWord := c.Request.URL.Query().Get("password")
	rows,err := DB.Query("SELECT * FROM voro_user") //查询用户名是否已存在
	if err != nil {
		fmt.Println("查询失败")
	}
	for rows.Next(){
		var s User
		err = rows.Scan(&s.id,&s.username,&s.password)
		if err != nil{
			fmt.Println(err)
		}
		fmt.Println(s.username)
		if userName != s.username{
			// 执行插入
			result, err := DB.Exec("INSERT INTO voro_user(username,password,tel)VALUES (?,?,?)",userName,passWord,userTel)
			if err != nil {
				fmt.Println("执行失败")
				return
			} else {
				rows,_ := result.RowsAffected() //输出执行的行数
				if rows != 1{
					c.JSON(200,gin.H{
						"success":false,
					})
				} else {
					c.JSON(200,gin.H{ //注册成功
						"success":true,
						"username":userName,
					})
				}
			}
		} else {
			fmt.Println("用户名已被注册") //用户名已存在(注册失败 用户名已被注册)
			c.JSON(200,gin.H{
				"code":400,
				"success":false,
				"msg":"用户名已被注册",
			})
		}
	}
	rows.Close()
}

func main(){
	initDB()
	router := gin.Default()
	user := router.Group("/user")
	{
		user.POST("/login",userLogin)
		user.POST("/register",userRegister)
		user.POST("/forgetpassword",forgetPassword)
	}
	router.Run(":9000")
}
