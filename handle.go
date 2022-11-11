package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	pb "github.com/Sotatek-ThomasNguyen2/example-service.git/proto"
	"github.com/Sotatek-ThomasNguyen2/example-service.git/utils"
)

func (r *Router) JSON(code int, payload interface{}, ctx *gin.Context) {
	now := ctx.MustGet("now")
	if now != nil {
		s := time.Since(now.(time.Time))
		ctx.Header("x-api-duration", s.String())
		ctx.JSON(code, payload)
		return
	}
	ctx.JSON(code, payload)
}

func (r *Router) handleListUserPartner(ctx *gin.Context) {
	rq := &pb.UserPartnerRequest{}
	err := utils.BindQuery(rq, ctx)
	if err != nil {
		log.Println("err: ", err)
	}
	log.Println("r.u: ", r.u)
	ups, err := r.u.ListUserPartners(ctx, rq)
	log.Println("ups: ", ups)
	if err != nil {
		r.JSON(200, gin.H{"code": -1, "message": "not found user partner"}, ctx)
		return
	}
	r.JSON(200, ups, ctx)
}
