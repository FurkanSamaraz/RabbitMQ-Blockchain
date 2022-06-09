package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/streadway/amqp"
)

func main() {
	// Create a new RabbitMQ connection.
	connRabbitMQ, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	// Create a new Fiber instance.
	app := fiber.New()

	// Add middleware.
	app.Use(
		logger.New(), // add simple logger
	)

	// Add route.
	app.Get("/send", func(c *fiber.Ctx) error {
		// Checking, if query is empty.
		if c.Query("msg") == "" {
			log.Println("Missing 'msg' query parameter")
		}

		// Let's start by opening a channel to our RabbitMQ instance
		// over the connection we have already established
		ch, err := connRabbitMQ.Channel()
		if err != nil {
			return err
		}
		defer ch.Close()

		// With this channel open, we can then start to interact.
		// With the instance and declare Queues that we can publish and subscribe to.
		kuyruk, err := ch.QueueDeclare(
			"kuyruk1", // kuyruğumuzun ismi
			false,     // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,
		)
		// Handle any errors if we were unable to create the queue.
		if err != nil {
			return err
		}

		// Attempt to publish a message to the queue.
		err = ch.Publish(
			"",          // exchange
			kuyruk.Name, // Gönderilecek kuyruk ismi. Bu şekilde önceki oluşturduğumuz kuyruğun ismini alabiliriz
			false,       // mandatory
			false,       // immediate
			amqp.Publishing{
				ContentType: "text/plain", //mesajımızın tipi
				Body:        []byte("mesajımız"),
			},
		)
		if err != nil {
			return err
		}

		return nil
	})

	// Start Fiber API server.
	log.Fatal(app.Listen(":3000"))
}

/*package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Employee struct {
	gorm.Model
	FirstName string `gorm:"not null;size:30"`
	LastName  string `gorm:"not null;size:30"`
	Emails    []Email
}

type Email struct {
	gorm.Model
	EmployeeID int    `gorm:"index"`
	Mail       string `gorm:"type:varchar(50);unique_index"`
	IsActive   bool
}

func main() {
	db, err := gorm.Open("sqlite3", "db\\southwind.sdb")
	db.LogMode(true)
	defer db.Close()
	if err == nil {
		//db.SingularTable(true)
		db.AutoMigrate(&Employee{}, &Email{})
		db.Model(&Employee{}).Related(&Email{})

		burakMails := []Email{
			Email{Mail: "selim@buraksenyurt.com", IsActive: true},
			Email{Mail: "burak.senyurt@southwind.com", IsActive: false},
			Email{Mail: "burakselimsenyurt@gmail.com", IsActive: true},
		}

		burak := Employee{FirstName: "burak", LastName: "senyurt", Emails: burakMails}
		db.Create(&burak)

		loraMails := []Email{
			Email{Mail: "lora@kimbilll.moon", IsActive: true},
			Email{Mail: "kimbill.the.black.lora@southwind.com", IsActive: true},
		}
		lora := Employee{FirstName: "Lora", LastName: "Kimbılll", Emails: loraMails}
		db.Create(&lora)

		WriteToScreen(burak)
		WriteToScreen(lora)

		var burki Employee
		db.Find(&burki, "ID=?", 1) //Önce
		db.Model(&burki).Update("LastName", "Selim Senyurt")
		WriteToScreen(burki)

		var buffon Employee

		db.Model(&buffon).Where("ID=?", 2).Updates(map[string]interface{}{"FirstName": "Cianluici", "LastName": "Buffon"})
		db.First(&buffon, 2) //Direkt primary key üstünden(varsayılan olarak ID) arama yapar
		WriteToScreen(buffon)
	} else {
		fmt.Println(err.Error())
	}
}

func WriteToScreen(e Employee) {
	fmt.Printf("%d\t%s,%s,%s\n", e.ID, e.FirstName, e.LastName, e.CreatedAt)
	for _, email := range e.Emails {
		fmt.Printf("\t%d:%s\n", email.ID, email.Mail)
	}
}
*/
