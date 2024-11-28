package v1main

import (
	"fyne.io/fyne/v2"
	"github.com/SuperH-0630/hdangan/src/runtime"
	"math/rand"
	"time"
)

type Sentence struct {
	Title   string
	Message string
}

var rander = rand.New(rand.NewSource(time.Now().UnixNano()))

var SS = []Sentence{
	{
		Title:   "小王子说",
		Message: "大人热爱数字，如果你跟他们说你认识了新朋友，他们从来不会问你重要的事情。",
	},
	{
		Title:   "小王子说",
		Message: "所有大人最初都是孩子，但这很少有人记得。",
	},
	{
		Title:   "小王子说",
		Message: "大人自己什么都不懂，总是要小孩来给他们解释，这让我觉得很累。",
	},
	{
		Title:   "小王子说",
		Message: "就算他不停地向前走，那也走不了多远。",
	},
	{
		Title:   "小王子说",
		Message: "忘记朋友是很可悲的，不是每个人都有过朋友。",
	},

	{
		Title:   "小王子说",
		Message: "但可惜的是，我不能看见箱子里的绵羊。也许是我有点像大人了，我肯定已经变老了。",
	},
	{
		Title:   "小王子说",
		Message: "猴面包树在长大之前也是小树苗啊。",
	},
	{
		Title:   "小王子说",
		Message: "如果星球太小，而猴面包树又太多的话，星球最后将会被撑得爆裂。",
	},
	{
		Title:   "小王子说",
		Message: "原来在很长的时间里，你唯一的消遣是默默地欣赏日落。。",
	},
	{
		Title:   "小王子说",
		Message: "你知道吗？人在难过的时候就会爱上日落。",
	},

	{
		Title:   "小王子说",
		Message: "花是弱小的、淳樸的，它們總是設法保護自己，以為有了刺就可以顯出自己的厲害",
	},
	{
		Title:   "小王子说",
		Message: "如果一個人在幾百萬顆星星當中，愛上獨一無二的一朵花，那麼他只要望著，就會很快樂",
	},
	{
		Title:   "小王子说",
		Message: "我當時什麼也不懂！我應該根據她的行動，而不是根據她的話判斷她。我本應看出她耍那些手段背後的溫情，我當時太年青，還不懂得愛她。",
	},
	{
		Title:   "小王子说",
		Message: "因為你為你的玫瑰花了那麼多時間，它才變得那麼重要。",
	},
	{
		Title:   "狐狸说",
		Message: "你再去看看那些玫瑰吧。你會明白，你的玫瑰確實是世界上獨一無二的。",
	},
	{
		Title:   "小王子回去看那些玫瑰",
		Message: "你們很美，但是你們空虛。沒有人會為你們而死。一個普通的路人會覺得我的玫瑰和你們很像。但是，就她一朵，比你們全體更重要。",
	},
	{
		Title:   "狐狸说",
		Message: "如果你馴養了我，我們就彼此需要了。對我，你就是世界上獨一無二的；對你，我也是世界上獨一無二的。",
	},
	{
		Title:   "狐狸说",
		Message: "人類再沒有時間去了解其他事物了，他們總是到商店買現成的東西，不過世界上還沒有可以購買朋友的商店。",
	},
	{
		Title:   "狐狸說",
		Message: "這就是我的秘密。它很簡單：只有用心看，才能看清楚。重要的東西是眼睛看不見的。",
	},

	{
		Title:   "那时",
		Message: "它只是與成千上萬隻狐狸一樣的一隻狐狸。但是，我和它做了朋友，現在它就是世界上獨一無二的了",
	},

	{
		Title:   "小王子",
		Message: "一個人賣弄聰明，免不了會說謊話。",
	},

	{
		Title:   "小王子",
		Message: "你晚上仰望天空時，因為我住在其中的一顆星星上，因為我會在其中的一顆星星上笑，你會覺得所有的星星都在笑。",
	},
}

const key = "game"

const (
	Nothing = iota
	TurnOn
	TurnOff
)

var status = TurnOff

func StartTheGame(rt runtime.RunTime) {
	p := rt.App().Preferences()
	switch p.Int(key) {
	default:
		p.SetInt(key, TurnOn)
		TurnOnLucky(rt)
	case Nothing, TurnOn:
		p.SetInt(key, TurnOn)
		TurnOnLucky(rt)
	case TurnOff:
		TurnOffLucky(rt)
	}
}

func ChangeGame(rt runtime.RunTime) int {
	p := rt.App().Preferences()
	switch status {
	default:
		p.SetInt(key, TurnOn)
		TurnOnLucky(rt)
		return TurnOff
	case TurnOff:
		p.SetInt(key, TurnOn)
		TurnOnLucky(rt)
		return TurnOn
	case TurnOn:
		p.SetInt(key, TurnOff)
		TurnOffLucky(rt)
		return TurnOff
	}
}

func TurnOnLucky(rt runtime.RunTime) {
	if status == TurnOff {
		Game(rt)
	}
	return
}

func TurnOffLucky(rt runtime.RunTime) {
	if status == TurnOn {
		rt.StopGame()
	}
	return
}

func Game(rt runtime.RunTime) func() {
	tt := timeTicker(15 * time.Minute)
	stop := make(chan bool, 1)

	go func(rt runtime.RunTime) {
		defer func() {
			close(stop)
			close(tt)
		}()

		status = TurnOn
		for {
			select {
			case <-stop:
				status = TurnOff
				return
			case <-tt:
				showLucy(rt)
			}
		}
	}(rt)

	stopFunc := func() {
		stop <- true
	}

	rt.SetGameStopFunc(stopFunc)

	return stopFunc
}

func timeTicker(t time.Duration) chan bool {
	res := make(chan bool, 10)
	res <- true
	go func() {
		for range time.Tick(t) {
			res <- true
		}
	}()
	return res
}

func showLucy(rt runtime.RunTime) {
	noti := getNoti()
	rt.App().SendNotification(noti)
}

func getNoti() *fyne.Notification {
	i := rander.Intn(len(SS))
	s := SS[i]
	return fyne.NewNotification(s.Title, s.Message)
}
