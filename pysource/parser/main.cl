; UTILS
(def! getaddr (fn* [data index]
		   { :line 1 :col 2 }
))

(def! char (fn* [cha]
	(fn* [state]
		(if (= (count (:data state)) 0)
			{ :data "" :ast (:ast state) }
			(let* [c (subs (:data state) 0 1)]
				(if (= c cha)
				  {
					:data (subs (:data state) 1)
					:ast (str (:ast state) c)
				  }
				  state))))))

(def! choose (fn* [& args]
	(fn* [state]
		(if (= (count args) 0)
			{ :valid false :message (str "Could not find matching rules in: " args)}

			(let* [res (apply (first args) state)]
				(if (= res state)
					((apply choose (rest args)) state)
					res))))))

(def! repeat (fn* [fnc arg]
	(fn* [state]
		(let* [res (arg state)]
			(if (= res state)
				res
				{ :data (fnc (:data res) ()) :ast () })))))

; SYNTAX

(def! nums (fn* [state]
	(if (= (count (:data state)) 0)
		{ :data "" :ast (:ast state) }

		(let* [c (ord (subs (:data state) 0 1))]
			(if (&& (<= 48 c) (<= c 57))
				(nums {
					:data (subs (:data state) 1)
					:ast (str (:ast state) (chr c))
				})
				state)))))

(def! ALPHA (fn* [state]
	(if (= (count (:data state)) 0)
		{ :data "" :ast (:ast state) }

		(let* [c (ord (subs (:data state) 0 1))]
			(if (&& (<= 65 c) (<= c 90))
				(nums {
					:data (subs (:data state) 1)
					:ast (str (:ast state) (chr c))
				})
				state)))))

(def! alpha (fn* [state]
	(if (= (count (:data state)) 0)
		{ :data "" :ast (:ast state) }

		(let* [c (ord (subs (:data state) 0 1))]
			(if (&& (<= 97 c) (<= c 122))
				(nums {
					:data (subs (:data state) 1)
					:ast (str (:ast state) (chr c))
				})
				state)))))

(def! whitespace (fn* [state]
	(choose ;comment
		(char "\n")
	  	(char "\t")
		(char " ")
		(char ","))))

(def! blist (fn* []
	(sequence (char "(")
		  (repeat bexpr)
		  (char ")"))))

(def! bkeyword (fn* []
	(sequence (char ":")
		  (choose alpha
			  ALPHA
			  (char "_") 
			  (char "-"))
		  (repeat (choose alpha
				  ALPHA
				  nums
				  (char "_") 
				  (char "-"))))))

(def! bexpr (fn* []
	(choose blist
		bkeyword)))

(def! test (repeat list (char "6")))
(def! test2 (repeat str (char "6")))

; ENV
(def! data "666666969antonin (1 2 4 \"lol\" nil true)")
(def! state { :data data :ast "" })
(prn state)
(prn (test state))
(prn (test2 state))



