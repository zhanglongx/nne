package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gotk3/gotk3/gtk"
	"github.com/zhanglongx/nns/kq"
)

var version string = "0.9.0 (Beta)"

var workers map[int]string = map[int]string{
	11613: "董家炜",
	11898: "赵晨辉",
	11956: "庞志远",
	11340: "魏永彬",
	13059: "韩飞",
	40531: "刘灿",
	12845: "李学良",
	13057: "李树超",
	13185: "姜东超",
	12585: "冯凯凯",
	40173: "李新",
	13101: "王炳建",
	13055: "马辉辉",
	23130: "闫小超",

	11479: "张龙霄",
}

func main() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("KQ  " + version)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Create a new grid widget to arrange child widgets
	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}

	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	lab1, err := gtk.LabelNew("输入起始日期, eg. 2006-01-02")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	entry1, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry:", err)
	}

	lab2, err := gtk.LabelNew("输入截止日期, eg. 2006-01-31")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	entry2, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry:", err)
	}

	lab3, err := gtk.LabelNew("输入存放目录, eg. D:\\")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	entry3, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create entry:", err)
	}

	btn, err := gtk.ButtonNewWithLabel("统计")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}

	lab4, err := gtk.LabelNew("作者：")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	lab5, err := gtk.LabelNew("zhanglongx@gmail.com under GPL-v3")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	btn.Connect("clicked", func() {
		start, _ := entry1.GetText()
		end, _ := entry2.GetText()

		file, _ := entry3.GetText()

		err := run(start, end, file)
		if err != nil {
			msg := gtk.MessageDialogNew(nil, gtk.DIALOG_MODAL,
				gtk.MESSAGE_ERROR,
				gtk.BUTTONS_NONE,
				fmt.Sprintf("%s", err),
			)

			msg.SetTitle("出错啦")

			msg.Show()
		}
	})

	grid.Add(lab1)
	grid.Add(entry1)
	grid.Add(lab2)
	grid.Add(entry2)
	grid.Add(lab3)
	grid.Add(entry3)
	grid.Add(btn)
	grid.Add(lab4)
	grid.Add(lab5)

	// Add the grid to the window, and show all widgets.
	win.Add(grid)

	// Set the default window size.
	win.SetDefaultSize(300, 300)

	win.ShowAll()

	gtk.Main()
}

func run(s string, e string, file string) error {

	start, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	end, err := time.Parse("2006-01-02", e)
	if err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close()

	ids := make([]int, len(workers))
	i := 0
	for id := range workers {
		ids[i] = id
		i++
	}

	k := kq.KQ{
		StartDate: start,
		EndDate:   end,

		Workers: ids,
	}

	info, err := k.Run()
	if err != nil {
		return err
	}

	f.WriteString(fmt.Sprintf("工号, 姓名, 补贴1, 补贴2\n"))
	for id, row := range info {
		f.WriteString(fmt.Sprintf("%d, %s, %d, %d\n",
			id, workers[id], row.Days1, row.Days2))
	}

	return nil
}
