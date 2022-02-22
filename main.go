package main

type User interface {
	Kick()
	Stop()
}

type Admin struct {
	name string
}

func (this Admin) Kick() {
	println(this.name + " KICK")
}

func (this Admin) Stop() {
	println("STOP")
}

func play(user User) {
	user.Kick()
}

func main() {
	a := []User{Admin{"Edgard"}, Admin{"Pedro"}}
	for _, i := range a {
		i.Kick()
	}
}
