package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type BlockChain struct {
	blocks []*Block
}

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

func (b *Block) DeriveHash() {
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(info)
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash}
	block.DeriveHash()
	return block
}

func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.blocks[len(chain.blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.blocks = append(chain.blocks, new)
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}
func main() {
	//RabbitMQ Sunucumuza bağlanıyoruz
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	defer conn.Close()

	//İletişim kurabilmek için kanal oluşturalım
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln(err)
	}
	defer ch.Close()

	//Kuyruğumuzu tanımlıyoruz
	_, err = ch.QueueDeclare(
		"kuyruk1", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalln(err)
	}

	//İşte burada kuyruğumuzu dinliyoruz.
	msgs, err := ch.Consume(
		"kuyruk1", // Bu sfer dinleyeceğim kuyruk ismini kendim yazdım
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalln(err)
	}
	//Burada goroutine ile çalışan fonksiyonumuz
	//çalışırken programın kapanmaması için
	//kanal oluşturduk
	forever := make(chan bool)
	chain := InitBlockChain()
	go func() {
		//Burada eğer varsa kuyruktaki mesajları çekiyoruz
		for d := range msgs {
			//d değişkeni ile kuyruktaki mesajın bilgilerine ulaşabiliriz.
			log.Printf("Alınan mesaj: %s", d.Body)
			//Kuyruktaki mesaj ekrana bastırdık.
			//kuyruktan alınan mesajı blocklama işlemine soktuk.
			chain.AddBlock(string(d.Body))

			for _, block := range chain.blocks {
				fmt.Printf("Previous Hash: %x\n", block.PrevHash)
				fmt.Printf("Data in Block: %s\n", block.Data)
				fmt.Printf("Hash: %x\n", block.Hash)
			}
		}
	}()

	log.Printf(" [*] Kuyruk1 dinleniyor...")

	//Burada forever isimli kanalımıza değer gönderilmeyeceği için
	//programımız kapanmayacak ve sürekli olarak kuyruktaki mesajları çekecektir.
	<-forever
}
