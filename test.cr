(import io)
(import arr)

"C'est un commentaire en texte"
(defunc greeting (name)
  (+ "Hello " name " !")
)

(defunc greet_everyones_except (names except)  ; qdsfqsd
  (->
    names
    (filter (func (name) (!= name except)))
    (map greeting)
    (join "\n")
  )
)

(io.echo 1)

(let group (array "Jackie" "Da\n\\\"niel" "Jean" "Paul"))

(io.echo (greet_everyones_except group "Jean"))
