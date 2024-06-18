package main

import (
	"fmt"
	"math/rand"
	"strings"
)

var prefixes = []string{
	"A painting",
	"A watercolour painting",
	"An impressionist painting",
	"A photographic portrait",
	"A drawing",
	"A sketch",
	"An illustration",
	"A photograph",
	"A digital sketch",
	"A photorealistic picture",
	"A computer render",
}

var artists = []string{
	"Albrecht Durer",
	"Alex Grey",
	"Alex Gross",
	"Alice Neel",
	"Amedeo Modigliani",
	"Andre Kohn",
	"Andrea Kowch",
	"Andrei Rublev",
	"Andrew Atroshenko",
	"Andrew Macara",
	"Andrey Remnev",
	"Andy Warhol",
	"Anne Stokes",
	"Anthony van Dyck",
	"Antonello da Messina",
	"August Sander",
	"A-1 Pictures",
	"Banksy",
	"Bella Kotaki",
	"Berthe Morisot",
	"Bill Gekas",
	"Bob Peak",
	"Botero",
	"Brad Kunkle",
	"Brian De Palma",
	"Carl Larsson",
	"Cartoon Network",
	"Charlie Bowater",
	"Childe Hassam",
	"Chuck Close",
	"Cimabue",
	"Cindy Sherman",
	"Craig Davison",
	"Dante Gabriel Rossetti",
	"Dario Argento",
	"Diego Rivera",
	"Donatello",
	"Edith Head",
	"Eduardo Kobra",
	"Edvard Munch",
	"Eric Wallis",
	"Ernst Ludwig Kirchner",
	"Francis Bacon",
	"Georges Seurat",
	"Gilbert Stuart",
	"Hans Holbein the Elder",
	"Hideaki Anno",
	"Ilya Kuvshinov",
	"Jean-Baptiste Monge",
	"Jim Henson",
	"Joaquín Sorolla",
	"Jules Bastien-Lepage",
	"Ken Kelly",
	"Krenz Cushart",
	"Leonardo da Vinci",
	"Lucian Freud",
	"Luis Royo",
	"Mark Keathley",
	"Pablo Picasso",
	"Patrice Murciano",
	"Paul Gauguin",
	"Quentin Blake",
	"Raphael",
	"Rembrandt",
	"Rene Magritte",
	"RossDraws",
	"Russ Mills",
	"Salvador Dalí",
	"Sandro Botticelli",
	"Sharaku",
	"Stevan Dohanos",
	"Steve McCurry",
	"Thomas Gainsborough",
	"Tom Bagshaw",
	"Vincent van Gogh",
	"Willem de Kooning",
}

var films = []string{
	"Harry Potter",
	"The Matrix",
	"The Martian",
	"The Lord of the Rings",
	"Pulp Fiction",
	"Fight Club",
	"Star Wars",
	"Back to the Future",
	"Aliens",
	"2001 A Space Odyssey",
	"Toy Story",
	"Hamilton",
	"Die Hard",
	"Indiana Jones and the Last Crusade",
	"James Bond",
	"Jurassic Park",
	"Blade Runner",
	"The Simpsons",
	"Futurama",
	"Dune",
	"Star Trek The Next Generation",
	"Star Trek Voyager",
	"Star Trek Deep Space Nine",
	"The Lion King",
}

var suffixes = []string{
	"bokeh",
	"long exposure",
	"unreal engine",
	"octane rendering",
	"8K",
	"trending in Artstation",
	"stunning",
	"dramatic cinematic lighting",
	"smiling",
	"cyberpunk",
	"national geographic",
	"inspired by Neal Stephenson",
	"sci-fi theme",
	"masterpiece",
	"award winning",
	"superb",
	"outstanding",
	"photoshopped",
	"hyper realism",
	"highly detailed",
	"pencil and paper",
	"charcoal on paper",
	"trading card foil",
	"prize winning",
	"board game art",
	"video game screenshot",
	"cinestill",
	"back light",
	"studio light",
	"natural framing",
	"soft focus",
	"technicolor",
	"asymmetrical composition",
	"extreme tonal balance",
}

const subject = "cjs"

func Prompt() string {
	prefix := prefixes[rand.Intn(len(prefixes))]
	extraCount := 1 + rand.Intn(2) + rand.Intn(3) + rand.Intn(4)
	r := rand.Intn(100)
	useArtist := r > 40
	useFilm := r < 30

	var trailing []string

	if useArtist {
		artist := artists[rand.Intn(len(artists))]
		trailing = append(trailing, fmt.Sprintf("in the style of %s", artist))
	}

	if useFilm {
		film := films[rand.Intn(len(films))]
		trailing = append(trailing, fmt.Sprintf("from the film %s", film))
	}

	rand.Shuffle(len(suffixes), func(i, j int) {
		suffixes[i], suffixes[j] = suffixes[j], suffixes[i]
	})

	trailing = append(trailing, suffixes[:extraCount]...)

	return fmt.Sprintf("%s of (((%s))), %s", prefix, subject, strings.Join(trailing, ", "))
}
