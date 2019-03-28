package corll

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"log"
	"myProject/db"
	"net/http"
	"strconv"
)

/*login user*/
func Login(g *gin.Context) {
	fmt.Println("login.........")
	rsp := new(Rsp)
	name := g.PostForm("username")
	pass := g.PostForm("password")
	mgo := db.InitMongoDB()
	findfilter := bson.D{{"username", g.PostForm("username")}, {"password", g.PostForm("password")}}
	cur, err := mgo.Collection(db.User).Find(context.Background(), findfilter)
	for cur.Next(context.Background()) {
		elme := new(User)
		err := cur.Decode(elme)
		if err == nil {
			if elme.Username == name && elme.Password == pass {
				rsp.Msg = "success"
				rsp.Code = 200
				g.JSON(http.StatusOK, rsp)
				return
			}
		}
	}

	rsp.Msg = "user is null"
	rsp.Code = 201
	rsp.Data = err
	g.JSON(http.StatusOK, rsp)
}

/* insert user table */
func Insertuser(g *gin.Context) {
	rsp := new(Rsp)
	fmt.Println("InsertUser.........")
	if g.PostForm("username") == "" || g.PostForm("password") == "" {
		rsp.Msg = "user is null"
		rsp.Code = 201
		g.JSON(http.StatusOK, rsp)
	}
	mgo := db.InitMongoDB()
	newuser := new(User)
	newuser.Username = g.PostForm("username")
	newuser.Password = g.PostForm("password")
	insertID, err := mgo.Collection(db.User).InsertOne(context.Background(), newuser)
	fmt.Println(insertID)
	if err == nil {
		rsp.Msg = "success"
		rsp.Code = 200
		g.JSON(http.StatusOK, rsp)
		return
	} else {
		rsp.Msg = "faild"
		rsp.Code = 201
		g.JSON(http.StatusOK, rsp)
		return
	}
}

/* query all user */
func Queryalluser(g *gin.Context) {
	fmt.Println("Queryalluser.........")
	rsp := new(Rsp)
	mgo := db.InitMongoDB()
	var users []User
	cur, err := mgo.Collection(db.User).Find(context.Background(), bson.D{}, nil)
	if err == nil {
		for cur.Next(context.Background()) {
			elme := new(User)
			err := cur.Decode(elme)
			if err == nil {
				users = append(users, *elme)
			}
		}
	}
	rsp.Msg = "success"
	rsp.Code = 200
	rsp.Data = users
	g.JSON(http.StatusOK, rsp)
	return
}

/*  模糊查询 query  by username */
func QueryByUsername(g *gin.Context) {
	rsp := new(Rsp)
	fmt.Println(".....QueryByUsername..")
	username := g.Query("username")
	mgo := db.InitMongoDB()
	fmt.Println(username)

	filter := bson.M{"username": bson.M{"$regex": username, "$options": "$i"}}

	cur, err := mgo.Collection(db.User).Find(context.Background(), filter)
	if err != nil {
		fmt.Println(err)
	}
	var users []User
	for cur.Next(context.Background()) {
		elme := new(User)
		err := cur.Decode(elme)
		if err == nil {
			users = append(users, *elme)
		}
	}
	rsp.Msg = "success"
	rsp.Code = 200
	rsp.Data = users
	g.JSON(http.StatusOK, rsp)
	return
}

/* get all user */
func Getalluser(g *gin.Context) {
	log.Println("Getalluser.........")
	rsp := new(Rsp)
	mgo := db.InitMongoDB()
	filter := bson.M{}
	limit, err := strconv.Atoi(g.Query("limit"))

	//排序 正序1 倒序-1  ----------------------------
	opts := new(options.FindOptions)
	sortMap := make(map[string]interface{})
	sortMap["gender"] = -1
	opts.Sort = sortMap
	//排序 正序1 倒序-1  ----------------------------

	var users []User
	cur, err := mgo.Collection(db.User).Find(context.Background(), filter, opts.SetLimit(int64(limit)))
	if err == nil {
		for cur.Next(context.Background()) {
			elme := new(User)
			err := cur.Decode(elme)
			if err == nil {
				users = append(users, *elme)
			}
		}
	}

	rsp.Msg = "success"
	rsp.Code = 0
	rsp.Data = users
	g.JSON(http.StatusOK, rsp)
	return
}

/* update user */
func Updateuser(g *gin.Context) {
	fmt.Println("update user.................")
	rsp := new(Rsp)
	mgo := db.InitMongoDB()
	id := g.PostForm("id")
	oldId, err := primitive.ObjectIDFromHex(id)
	newpassword := g.PostForm("password")
	newgender, err := strconv.Atoi(g.PostForm("gender"))

	user := new(User)
	user.Id = oldId

	oldInfo := mgo.Collection(db.User).FindOne(context.Background(), bson.M{"_id": oldId})
	if oldInfo != nil {
		oldInfo.Decode(user)
	}
	if g.PostForm("username") != "" {
		user.Username = g.PostForm("username")
	}
	if newpassword != "" {
		user.Password = newpassword
	}
	if newgender != 0 {
		user.Gender = newgender
	}
	if g.PostForm("address") != "" {
		user.Address = g.PostForm("address")
	}

	//info, err := mgo.Collection(db.User).UpdateOne(context.Background(), bson.M{"_id": bsonx.ObjectID(oldId)},
	//	bson.M{"$set": bson.M{"username": newUsername}})
	//
	info := mgo.Collection(db.User).FindOneAndReplace(context.Background(), bson.M{"_id": oldId},
		user)
	if info.Err() != nil {
		rsp.Msg = info.Err().Error()
		rsp.Code = 201
		g.JSON(http.StatusOK, rsp)
		return
	}
	fmt.Println(info)
	if err == nil {
		//fmt.Println(info.MatchedCount)
	}
	rsp.Msg = "success"
	rsp.Code = 200
	g.JSON(http.StatusOK, rsp)
	return
}

/* deluser*/
func Deluser(g *gin.Context) {
	fmt.Println("update user.................")
	rsp := new(Rsp)
	mgo := db.InitMongoDB()
	username := g.PostForm("username")
	info, err := mgo.Collection(db.User).DeleteOne(context.Background(), bson.M{"username": username})
	if info.DeletedCount == 0 {
		rsp.Msg = "faild"
		rsp.Data = "username is  not exist!!"
		rsp.Code = 201
		g.JSON(http.StatusOK, rsp)
		return
	}
	if err == nil {
		fmt.Println(info.DeletedCount)
		rsp.Msg = "success"
		rsp.Code = 200
		g.JSON(http.StatusOK, rsp)
		return
	} else {
		rsp.Msg = "faild"
		rsp.Data = err
		rsp.Code = 201
		g.JSON(http.StatusOK, rsp)
		return
	}
}