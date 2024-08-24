package main

import (
	"fmt"
	"strings"
)

/*
	05.07
		ПРОВЕДИ ТЕСТИРОВАНИЕ НОРМАЛЬНОЕ, У ТЕБЯ КУЧА ОШИБОК
		ЧТО ОСТАЛОСЬ СДЕЛАТЬ:
			3. Добавить это сообщение на улице весна. можно пройти - домой
*/

/*
	код писать в этом файле
	наверняка у вас будут какие-то структуры с методами, глобальные переменные ( тут можно ), функции
*/

type Room struct {
	Name       string
	Furniture  map[string]*FurnitureItem
	Directions map[string]*Room
	Doors      map[string]*Door
	furnitureKeys []string
	DirectionsKeys []string
}

type Door struct {
	isClosed bool
}

type FurnitureItem struct {
	Name      string
	Inventory map[string]bool
	itemKeys []string
}

type Player struct {
	Enviroment string
	Bag        *Backpack
}

type Backpack struct {
	Inventory map[string]bool
	OnPlayer  bool
}


var rooms map[string]*Room = make(map[string]*Room)

var player Player

func main() {
	initGame()
	fmt.Println(handleCommand("осмотреться"))
	fmt.Println(handleCommand("идти коридор"))
	fmt.Println(handleCommand("идти комната"))
	fmt.Println(handleCommand("надеть рюкзак"))
	fmt.Println(handleCommand("взять ключи"))
	fmt.Println(handleCommand("идти коридор"))
	fmt.Println(handleCommand("применить ключи дверь"))
	
}

func initGame() {
	// Rooms := []*Room{}
	initPlayer()
	addRoom("кухня", &rooms)
	addRoom("коридор", &rooms)
	addRoom("комната", &rooms)
	addRoom("улица", &rooms)
	rooms["кухня"].AddPass(rooms["коридор"])
	rooms["комната"].AddPass(rooms["коридор"])
	rooms["улица"].AddPass(rooms["коридор"])
	rooms["коридор"].AddClosedDoor(rooms["улица"])
	fmt.Println(rooms["улица"].Doors["коридор"].isClosed)
	rooms["кухня"].AddFurniture("стол")
	rooms["кухня"].AddFurniture("стул")
	rooms["комната"].AddFurniture("стол")
	rooms["комната"].AddFurniture("стул")
	rooms["комната"].Furniture["стул"].AddItemFurniture("рюкзак") // можно добавить обработку ошибки при неправильном вводе комнаты
	rooms["комната"].Furniture["стол"].AddItemFurniture("ключи")
	rooms["комната"].Furniture["стол"].AddItemFurniture("конспекты")
	rooms["кухня"].Furniture["стол"].AddItemFurniture("чай")
}

func (room Room) returnInventoryStr() string { // эта функция должна пре
	var inventoryMessage string
	if room.HasOneNonEmptyFurniture() {
		for _, name := range room.furnitureKeys {
			if room.Furniture[name].AtLeastOneItemInInventory() {
				inventoryMessage += "на " + room.Furniture[name].Name + "е: "
				for _, itemName := range room.Furniture[name].itemKeys {
					if room.Furniture[name].Inventory[itemName] {
						inventoryMessage += itemName + ", "
					}
				}
			}
		}
		inventoryMessage = inventoryMessage[:len(inventoryMessage) - 2] + "."
	} else {
		inventoryMessage = "пустая комната."
	}
	return inventoryMessage
}

func (room Room) HasOneNonEmptyFurniture() bool {
	result := false
	for _, val := range room.Furniture {
		if val.AtLeastOneItemInInventory() {
			result = true
			break
		}
	}
	return result
} 

func (to FurnitureItem) AtLeastOneItemInInventory() bool{ // для 1) 1.
	result := false
	for _, val := range to.Inventory{
		if val {
			result = true
			break
		}
	}
	return result
}


func (player *Player) lookAround(isNotWalk bool) string {
	/*
		Функция должна возвращать строку - результат на команду осмотреться
		Для этого: result складывается из трех строк-сообщений: where, furniture и directions
		
	*/


	var result, where, furniture, directions string
	if isNotWalk {
		furniture = rooms[player.Enviroment].returnInventoryStr()
	}
	
	// Доделать
	directions = rooms[player.Enviroment].returnDirectionsStr()
	
	switch player.Enviroment {
	case "кухня":
		var goal string
		if player.Bag.OnPlayer && player.Bag.Inventory["конспекты"] && player.Bag.Inventory["ключи"]{
			goal = "надо идти в универ."
		} else {
			goal = "надо собрать рюкзак и идти в универ."
		}
		if isNotWalk {
			furniture = furniture[:len(furniture) - 1] + ", " + goal
			where = "ты находишься на кухне, "
		} else {
			furniture += "." 
			where = "кухня, ничего интересного"
		}
	case "комната":
		if !isNotWalk {
			where = "ты в своей комнате."
		}
	case "коридор":
		where = "ничего интересного."
	case "улица":
		where = "на улице весна."
		directions = " можно пройти - домой"
	default:
		where = ""
	}
	
	result = where + furniture + directions
	return result
}

func (room Room) returnDirectionsStr() string {
	var directionsMessage string
	directionsMessage += " можно пройти - "
	for _, val := range room.DirectionsKeys {
		directionsMessage += val + ", "
	}
	directionsMessage = directionsMessage[:len(directionsMessage) - 2]
	return directionsMessage
}

func initPlayer() {
	player = Player{
		Enviroment: "кухня",
		Bag: &Backpack{
			Inventory: make(map[string]bool),
			OnPlayer:  false,
		},
	}
}

func (room *Room) AddFurniture(name string) {
	room.Furniture[name] = &FurnitureItem{
		Name:      name,
		Inventory: make(map[string]bool),
	}
	room.furnitureKeys = append(room.furnitureKeys, name)
}

func (to *FurnitureItem) AddItemFurniture(item string) {
	if !to.Inventory[item] {
		to.Inventory[item] = true
	}
	to.itemKeys = append(to.itemKeys, item)

}

func addRoom(name string, where *map[string]*Room) {
	room := &Room{
		Name:       name,
		Furniture:  make(map[string]*FurnitureItem),
		Directions: make(map[string]*Room),
		Doors:      make(map[string]*Door),
	}
	(*where)[name] = room
}

func (room *Room) AddClosedDoor(to *Room) { // возможно стоит добавить параметры замка и открыт закрыт тогда нужна будет отдельная структура Дверь с полями есть замок\нет открыта\закрыта
	if room.Directions[to.Name] == to { // есть подозрение что условие не отрабатывает
		// fmt.Println("OTRABOTAL")
		room.Doors[to.Name] = &Door{
			isClosed: true,
		}
		to.Doors[room.Name] = room.Doors[to.Name]
	}
	// fmt.Println(from.Doors[to.Name].isClosed)
	// fmt.Println(to.Doors[from.Name].isClosed)

}

func (room *Room) AddPass(to *Room) {
	room.Directions[to.Name] = to
	room.DirectionsKeys = append(room.DirectionsKeys, to.Name)
	to.Directions[room.Name] = room
	to.DirectionsKeys = append(to.DirectionsKeys, room.Name)
}

func handleCommand(command string) string {
	var message string
	commands := strings.Split(command, " ")
	switch commands[0] {
	case "осмотреться":
		message = player.lookAround(true)
	case "идти":
		where := commands[1]
		/*
			прежде чем пытаться получить доступ к двери, нужно проверить существует ли онаю
			Дверь существует, если указатель на нее не равень nil.
			Если дверь существует, то только тогда мы пытаемся получить доступ к isClosed
		*/
		if rooms[player.Enviroment].Directions[where] != nil {
			if rooms[player.Enviroment].Doors[where] != nil {
				if rooms[player.Enviroment].Doors[where].isClosed {
					message = "дверь закрыта"
				} else {
					player.Enviroment = where
					message = player.lookAround(false)
				}
			} else {
				player.Enviroment = where
				message = player.lookAround(false)
			}
		// } else if rooms[player.Enviroment].Doors[where].isClosed {
		// 	message = "дверь закрыта"
		} else {
			message = "нет пути в " + where
		}
	case "применить":
		what := commands[1]
		to := commands[2]
		message = rooms[player.Enviroment].Apply(what, to)
	case "взять":
		item := commands[1]
		message = rooms[player.Enviroment].Grab(item)
	case "надеть":
		message = rooms[player.Enviroment].PutOnBackpack()
	default: 
		message = "неизвестная команда"
	}
	return message
}

func (room *Room) PutOnBackpack() string {
	var result string
	itemFounded := false
	Loop:
	for _, val := range room.Furniture {
		for item, inStock := range val.Inventory {
			if item == "рюкзак" && inStock {
				player.Bag.OnPlayer = true
				val.Inventory[item] = false
				itemFounded = true
				result = "вы надели: рюкзак"
				break Loop
			}
		}
	}
	if !itemFounded {
		result = "нет такого"
	}
	return result
}

func (room *Room) Grab(something string) string{
	var result string
	/*
		нужна функция которая пробегается по всей мебели в комнате и ищет конкретный предмет. Возвращает  - если найдено, false - если нет
	*/
	if player.Bag.OnPlayer {
		itemFounded := false
		Loop:
		for _, val := range room.Furniture {
			for item, inStock := range val.Inventory {
				if item == something && inStock {
					player.Bag.Inventory[something] = true
					val.Inventory[item] = false
					itemFounded = true
					result = "предмет добавлен в инвентарь: " + something
					break Loop
				}
			}
		}
		if !itemFounded {
			result = "нет такого"
		}
	} else {
		result = "некуда класть"
	}
	return result
}



func (room *Room) Apply(what string, to string) string {
	var result string
	if player.Bag.Inventory[what] && player.Bag.OnPlayer {
		if room.Doors["улица"] != nil {
			if room.Doors["улица"].isClosed && what == "ключи" && to == "дверь" { // почему-то isClosed не существует
				result = "дверь открыта"
				room.Doors["улица"].isClosed = false
			} else {
				result = "не к чему применить"
			}
		} else {
			result = "не к чему применить"
		}
	} else {
		result = "нет предмета в инвентаре - " + what
	}
	return result
}