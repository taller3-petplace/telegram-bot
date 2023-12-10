package salchifacts

import (
	"fmt"
	"math/rand"
	"telegram-bot/internal/utils/formatter"
)

var facts = []string{
	origin,
	excavators,
	dna,
	popularity,
	badClimbers,
	war,
	energy,
	taylorSwim,
	hunters,
	bold,
	smell,
	smart,
	perseverance,
	song,
	formatter.Link("Watch this Reel", "https://www.instagram.com/reel/C0FEkDNgJhP/?igshid=MzRlODBiNWFlZA=="),
}

func GetFact() string {
	factNumber := rand.Intn(len(facts))
	fact := facts[factNumber]
	return fmt.Sprintf(
		"Salchi Fact #%v: \n%s",
		factNumber+1,
		fact,
	)
}

var origin = "El origen de esta raza se remonta a la Alemania del siglo XVII. Los perros salchichas eran usados, " +
	"principalmente, para cazar tejones. También era habitual que cazaran animales de madriguera como conejos y zorros"

var excavators = fmt.Sprintf("%s Su cuerpo alargado, las piernas cortas y fuertes, y sus afiladas garras los "+
	"convierten en excelentes excavadores y cazadores subterráneos.",
	formatter.Bold("Son excelentes excavadores:"),
)

var dna = "El aspecto tan particular se debe a una mutación genética conocida como bassetismo. " +
	"Por eso el nombre de los basset hound. Debido a esta condición, su columna es alargada y sus patas cortas."

var popularity = "Los perros salchichas son el símbolo nacional de Alemania. Además, es una raza muy popular. " +
	"El dachshund ocupa el puesto 13 entre las 194 razas de perros reconocidas por el American Kennel Club."

var badClimbers = "Por su peculiar forma, los perros salchicha son muy malos escaladores, " +
	"**por lo que no deben subir ni bajar escaleras**. Además, suelen sufrir problemas de columna."

var war = "Al ser tan estrechamente relacionados con Alemania, estos perros no fueron bien vistos durante las guerras mundiales."

var energy = "Una de sus características principales es su gran reserva de energía. Al ser unos perros  para la caza, aman estar al aire libre." +
	"\n\nPor ello requieren bastante actividad física, salir a trotar o correr al menos una hora al día. Más allá de eso tienen una actitud tranquila, " +
	"aunque son algo obstinados."

var taylorSwim = "Tienen una particular forma de nadar, por lo que a los pocos minutos tienden  irse de lado y tendrás que rescatarlos."

var hunters = "Conserva su instinto cazador, por lo que debemos estar atentos cuando lo soltemos. " +
	"Eso sí, con una correcta educación no tendría por qué desobedecer a nuestras órdenes."

var bold = "Es un perro bien equilibrado, con un porte audaz y seguro, con  una expresión facial inteligente y alerta."

var smell = "Los salchicha tienen un olfato extraordinario, bien constituido no podrás engañarle ya que tiene " +
	"espíritu de caza. Esto es debido a su pasado de cazador bajo tierra."

var smart = "Requiere cariño y compañía y es una raza muy sociable. " +
	"Además, a los salchicha les gusta jugar y tienen una rápida capacidad de aprendizaje"

var perseverance = "Debido a origen de cazadores han sido criados para ser perseverantes, " +
	"por lo que pueden llegar a ser algo tercos en el día a día."

var song = fmt.Sprintf(
	"Tienen una canción re piola, podes escucharla %s",
	formatter.Link("acá", "https://www.youtube.com/watch?v=0QAEwP9HO7M&ab_channel=Sujes"),
)
