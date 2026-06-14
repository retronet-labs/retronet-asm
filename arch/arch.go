// Package arch definisce il contratto che ogni architettura (4004, 8008, 8080)
// deve soddisfare per essere assemblata, e l'AST minimo condiviso.
package arch

// Instruction è una singola istruzione del sorgente: il mnemonico e i suoi
// operandi grezzi (ancora stringhe), più la riga per i messaggi d'errore.
type Instruction struct {
	Mnemonic string   // mnemonico normalizzato in MAIUSCOLO, es. "LDM"
	Operands []string // operandi grezzi, es. ["R4", "loop"]
	Line     int      // riga del sorgente (1-based), per gli errori
}

// Resolver restituisce l'indirizzo di una label e true se è definita.
type Resolver func(name string) (addr int, ok bool)

// Arch è il backend di una specifica architettura. L'assembler lo interroga in
// due passate: prima Size (per assegnare gli indirizzi alle label), poi Encode
// (per emettere i byte, ora che le label hanno un indirizzo noto).
type Arch interface {
	// Name è l'identificatore dell'architettura, es. "i4004".
	Name() string

	// Size restituisce la lunghezza in byte dell'istruzione (1 o 2 per il 4004),
	// o un errore se il mnemonico è sconosciuto o gli operandi non tornano.
	Size(in Instruction) (int, error)

	// Encode produce i byte dell'istruzione posta all'indirizzo pc, usando
	// resolve per tradurre eventuali label in indirizzi.
	Encode(in Instruction, pc int, resolve Resolver) ([]byte, error)
}
