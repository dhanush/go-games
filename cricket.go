//Demonstrates a cricket game score using unbuffered channels.
// Receives the input for the no of overs to be played.
//Generates a random score from 0-6 on each ball,
//and if the score is 5 then the batsman is out.
//Inspired by the unbuffered channel talk given by William Kennedy at Gophercon India 2015
package main

import (
	"fmt"
	"math/rand"
	"time"
)

//holds the score
type score struct {
	runs  int
	wkts  int
	overs int
}

//the bowler method just bowls the ball no received in the channel and calls the batter routine
//if the ball number is more than 6 we send a dummy value 100 which is discarded by the next over
func bowler(ch chan int, score *score) {
	fmt.Println("--------------------")
	i := <-ch

	if i > 6 {
		score.overs = score.overs + 1
		fmt.Printf("Over %d Complete \n", score.overs)
		fmt.Printf("Score at end of Over %d is %d/%d \n", score.overs, score.runs, score.wkts)
		ch <- 100
		return
	}
	fmt.Printf("Ball no %d bowled \n", i)
	go batter(ch, score)
	ch <- i
	time.Sleep(100 * time.Millisecond)
}

//batter hits the ball no received in the channel.
//Here he can get out too. In either case the ball no is incremented and the bowler routine is called
func batter(ch chan int, score *score) {
	i := <-ch
	//random no to calculate the runs it
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(100)
	runs := r % 6
	//if runs is 5 then it is considered out
	if runs == 5 {
		fmt.Printf("Ball no %d castles the batsman \n", i)
		score.wkts = score.wkts + 1
		//if all out then close the channel
		if score.wkts == 10 {
			fmt.Printf("All out. Final Score is %d/%d \n", score.runs, score.wkts)
			close(ch)
			return
		}
		fmt.Printf("Score now is %d/%d \n", score.runs, score.wkts)
		go bowler(ch, score)
		time.Sleep(100 * time.Millisecond)
		ch <- i + 1
		return
	}
	fmt.Printf("Ball no %d hit for %d runs by batsman \n", i, runs)
	score.runs = score.runs + runs
	fmt.Printf("Score now is %d/%d \n", score.runs, score.wkts)
	go bowler(ch, score)
	time.Sleep(100 * time.Millisecond)
	ch <- i + 1
}

func main() {

	ch := make(chan int)
	score := new(score)
	score.runs = 0
	score.wkts = 0
	score.overs = 0
	var o int
	fmt.Println("How many overs do you want to play ?")
	_, err := fmt.Scanf("%d", &o)
	if err != nil {
		panic(err)
	}
	for i := 1; i <= o; i++ {

		fmt.Println("######################")
		fmt.Printf("Over no %d starting \n", i)
		go bowler(ch, score)
		ch <- 1
		time.Sleep(2 * time.Second)
		//check for the allout condition. channel will be closed in that case and we break from the loop
		_, ok := <-ch
		if !ok {
			break
		}
	}
	fmt.Printf("Final Score is %d/%d in %d overs \n", score.runs, score.wkts, score.overs)
}
