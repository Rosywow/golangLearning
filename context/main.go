package main
//[go context解析 - 简书 (jianshu.com)](https://www.jianshu.com/p/737a115bc3cb)
//https://youtu.be/LSzR0VEraWw
import (
	"context"
	"fmt"
	"time"
)

func sleepAndTalk(ctx context.Context, duration time.Duration, s string) {

	select {
	case <-ctx.Done():
		fmt.Println("handle", ctx.Err())

	case <-time.After(duration):
		fmt.Println(s)
	}
}
var key ="k"
var apple = "a"
func main()  {
	ctx :=context.Background()
	//ctx,cancel :=context.WithCancel(ctx)
	ctx,cancel := context.WithTimeout(ctx,6*time.Second)

	defer cancel()
	//Value
	//ctx = context.WithValue(ctx,key,apple)
	//fmt.Println("key:",ctx.Value(key))


	//go func(){
	//	s:=bufio.NewScanner(os.Stdin)
	//	s.Scan()
	//	cancel()
	//}()

	sleepAndTalk(ctx,5*time.Second,"hello")
}
