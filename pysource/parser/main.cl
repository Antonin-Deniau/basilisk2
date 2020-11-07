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
				  state
				  ))))))

(def! choose (fn* [args]
	(fn* [state]
		(if (= (count args) 0)
			{ :valid false :message (prn "Could not find matching rules in: " args)}

			(let* [res ((get (first args)) state)]
				(if (true? (get res :valid))
				  res
				  (choose (rest args) state)))))))

(def! repeat (fn* []
	 1
))

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
	(choose (char "\n")
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

; ENV
(def! data "6969antonin (1 2 4 \"lol\" nil true)")
(def! state { :data data :ast "" })

(prn state)
(prn ((char "6") state))

