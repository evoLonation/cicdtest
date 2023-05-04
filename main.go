package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Res struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
type HelloReq struct {
	Content string `json:"content"`
}
type HelloRes struct {
	Msg string `json:"msg"`
}
type LoginReq struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
}
type LoginRes struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
type EnterItemRes struct {
	Price int `json:"price"`
}

// 该函数返回一个gin.H，gin.H是一个map，存储着键值对，将要返回给请求者
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

var uuidmp map[string]bool

func GetGinHandler[REQ any, RES any](logic func(req *REQ, id int) (*RES, error)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req REQ
		if err := ctx.ShouldBindJSON(&req); err != nil {
			//证明请求对于该结构体并不有效
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		studentId64, err := strconv.ParseInt(ctx.Request.Header.Get("StudentId"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, errorResponse(errors.New("请按照实验要求在请求头header中添加验证信息")))
			return
		}
		res, err := logic(&req, int(studentId64))
		uuid := uuid.New().String()
		uuidmp[uuid] = true
		log.Println(uuid)
		ctx.Header("Uuid", uuid)
		if err != nil {
			ctx.JSON(http.StatusPreconditionFailed, errorResponse(err))
		} else {
			ctx.JSON(http.StatusOK, res)
		}

	}
}

func GetGinHandlerNoReq[RES any](logic func(id int) (*RES, error)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		studentId64, err := strconv.ParseInt(ctx.Request.Header.Get("StudentId"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, errorResponse(errors.New("请按照实验要求在请求头header中添加验证信息")))
			return
		}
		res, err := logic(int(studentId64))
		uuid := uuid.New().String()
		uuidmp[uuid] = true
		log.Println(uuid)
		ctx.Header("UUID", uuid)
		if err != nil {
			ctx.JSON(http.StatusPreconditionFailed, errorResponse(err))
		} else {
			ctx.JSON(http.StatusOK, res)
		}

	}
}

func main() {
	uuidmp = make(map[string]bool)
	hellomp := make(map[int]bool)
	testmp := make(map[int]bool)
	frontmp := make(map[int]bool)
	router := gin.Default()
	router.LoadHTMLGlob("login.html")
	router.GET("zzy/login", GetGinHandler(func(req *LoginReq, id int) (*LoginRes, error) {
		if req.UserId == "zzy" {
			return nil, errors.New("giao")
		}
		return &LoginRes{Code: 0, Msg: "login success"}, nil
	}))
	router.POST("hello", GetGinHandler(func(req *HelloReq, id int) (*HelloRes, error) {
		if req.Content == "hello, world!" {
			hellomp[id] = true
			return &HelloRes{Msg: "hello, my friend!"}, nil
		} else {
			return nil, errors.New("请求体错误，检查content的值是否严格按照实验")
		}
	}))
	router.GET("hello/judge1", GetGinHandler(func(req *HelloReq, id int) (*Res, error) {
		if hellomp[id] {
			return &Res{Code: 0}, nil
		} else {
			return &Res{Code: 1, Msg: "请检查一下你的学号是否都输入正确，或者正确发送请求"}, nil
		}
	}))
	router.GET("hello/judge2", GetGinHandler(func(req *HelloReq, id int) (*Res, error) {
		if uuidmp[req.Content] {
			return &Res{Code: 0}, nil
		} else {
			return &Res{Code: 1, Msg: "请检查一下你的UUID是否输入正确"}, nil
		}
	}))
	router.GET("test", func(ctx *gin.Context) {

		studentId64, err := strconv.ParseInt(ctx.Request.Header.Get("StudentId"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, errorResponse(errors.New("请按照实验要求在请求头header中添加验证信息")))
			return
		}
		testmp[int(studentId64)] = true
		ctx.JSON(http.StatusPreconditionFailed, map[string]any{"code": 1})
	})
	router.GET("test/judge", GetGinHandler(func(req *HelloReq, id int) (*Res, error) {
		if testmp[id] {
			return &Res{Code: 0}, nil
		} else {
			return &Res{Code: 1, Msg: "请检查一下你的学号是否都输入正确，或者正确发送请求"}, nil
		}
	}))
	router.GET("front/login", func(ctx *gin.Context) {
		id := ctx.Query("id")
		studentId64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, errorResponse(errors.New("请按照实验要求在请求头header中添加验证信息")))
			return
		}
		frontmp[int(studentId64)] = true
		ctx.JSON(http.StatusOK, nil)
	})
	router.GET("login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", gin.H{})
	})
	router.GET("front/login/judge", GetGinHandlerNoReq(func(id int) (*Res, error) {
		if frontmp[id] {
			return &Res{Code: 0}, nil
		} else {
			return &Res{Code: 1, Msg: "请检查一下你的学号是否都输入正确，或者正确发送请求"}, nil
		}
	}))
	salemp := make(map[int]bool)
	pricemp := make(map[int]int)
	paymentmp := make(map[int]bool)
	router.POST("start-sale", GetGinHandlerNoReq(func(id int) (*Res, error) {
		salemp[id] = true
		return &Res{Code: 0}, nil
	}))
	router.POST("enter-item", GetGinHandler(func(req *HelloReq, id int) (*EnterItemRes, error) {
		if !salemp[id] {
			return nil, errors.New("请检查你是否事先执行startsale，或者设置了正确的学号")
		}
		pricemp[id] = pricemp[id] + 10
		return &EnterItemRes{Price: pricemp[id]}, nil
	}))
	router.POST("make-payment", GetGinHandler(func(req *EnterItemRes, id int) (*Res, error) {
		if !salemp[id] {
			return nil, errors.New("请检查你是否事先执行startsale，或者设置了正确的学号")
		}
		if req.Price != pricemp[id] {
			return nil, errors.New("支付金额不等于购买金额")
		}
		paymentmp[id] = true
		return &Res{Code: 0}, nil
	}))
	router.GET("process-sale/judge", GetGinHandlerNoReq(func(id int) (*Res, error) {
		if !paymentmp[id] {
			return nil, errors.New("请检查你是否完整执行了整个用例")
		}
		return &Res{Code: 0}, nil
	}))

	router.Run(":8081")
}
