package api

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 投票发起人在一个key对应的值里的第一个
// redis的键值是redis的topic加id，value值是[user][content]
type VoteRequest struct {
	VoteID   int           `json:"voteid"`
	Topic    string        `json:"topic"`
	Action   string        `json:"action"`
	Username string        `json:"username"`
	Content  string        `json:"Content"`
	Duration time.Duration `json:"duration"`
	Voting   Voting        `json:"voting"`
}

func (server *Server) vote(ctx *gin.Context) {
	var req VoteRequest
	var mu sync.Mutex
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, ctx.Request.Header)

	if err != nil {
		ctx.JSON(websocket.CloseProtocolError, errorResponse(err)) //1011表示websocket 服务器遇到异常
		conn.Close()
		return

	}
	if err = conn.ReadJSON(req); err != nil {
		ctx.JSON(websocket.CloseProtocolError, errorResponse(err)) //1011表示websocket 服务器遇到异常
		conn.Close()
		return
	}

	switch {

	//查看投票
	case req.Action == "check":
		result, err := server.cache.Get(ctx, req.Topic)
		if err != nil {
			ctx.JSON(websocket.CloseInternalServerErr, errorResponse(err))
			conn.Close()
			return
		}
		ctx.JSON(3001, normalResponce("check voting", result))
		return
	//发起投票
	case req.Action == "Initiate":
		mu.Lock()
		//这是第一个
		req.Voting.starter = req.Username
		req.Voting.VoteID = req.VoteID
		req.Voting.topic = req.Topic
		err := server.cache.Add(ctx, req.Topic, req.Voting, req.Duration)
		if err != nil {
			ctx.JSON(websocket.CloseInternalServerErr, errorResponse(err))
			mu.Unlock()
			conn.Close()
			return
		}
		mu.Unlock()
		ctx.JSON(3001, normalResponce("3001", "successfully initiate"))
		return
	//加入投票
	case req.Action == "join":
		mu.Lock()

		req.Voting.joiner[req.Username] = req.Content
		err := server.cache.Update(ctx, req.Topic, req.Voting, req.Duration)
		if err != nil {
			ctx.JSON(websocket.CloseInternalServerErr, errorResponse(err))
			mu.Unlock()
			conn.Close()
			return
		}
		mu.Unlock()
		ctx.JSON(3001, normalResponce("3001", "successfully join"))
		return
	case req.Action == "delete_topic":
		mu.Lock()

		req.Voting.joiner[req.Username] = req.Content
		err := server.cache.Delete(ctx, req.Topic)
		if err != nil {
			ctx.JSON(websocket.CloseInternalServerErr, errorResponse(err))
			mu.Unlock()
			conn.Close()
			return
		}
		mu.Unlock()
		ctx.JSON(3001, normalResponce("3001", "successfully delete topic"))
		return
	case req.Action == "delete_vote":
		mu.Lock()

		delete(req.Voting.joiner, req.Username)

		mu.Unlock()
		ctx.JSON(3001, normalResponce("3001", "successfully delete vote"))
		return
	default:
		ctx.JSON(1003, "connection closed")
		conn.Close()
		return
	}

}


