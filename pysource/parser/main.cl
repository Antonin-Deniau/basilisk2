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
			{ :valid false :message (prn "Could not find matching rules in: " args)}

			(let* [res (apply (first args) state)]
				(if (= res state)
					((apply choose (rest args)) state)
					res))))))

(def! repeat (fn* [typ fnc]
	(fn* [state]
		(let* [res (fnc state)]
			(if (= res state)
				state
				(let* [resp ((repeat typ fnc) { :data (:data res) :ast ""})]
				  { :data (:data resp) :ast (apply typ (cons (:ast res) (:ast resp))) })
				)))))

(def! sequence (fn* [& args]
	1
))

(def! ignore (fn* []
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

(def! whitespace (choose (char "\n")
			 (char "\t")
			 (char " ")
			 (char ",")))

;(def! blist (sequence (char "(")
;		      (repeat bexpr)
;		      (char ")")))

;(def! bexpr (choose blist
;		    bkeyword))


(def! bkeyword (sequence (char ":")
		         (choose alpha
				 ALPHA
				 (char "_") 
				 (char "-"))
		         (repeat vector (choose alpha
					        ALPHA
					        nums
					        (char "_") 
					        (char "-")))))

;(def! test (sequence (char "6")
;		     (char "9")
;		     (char "6")))
(def! test (repeat vector (char "6")))


; ENV
(def! data "666969antonin (1 2 4 \"lol\" nil true)")
(def! state { :data data :ast "" })
(prn state)
(prn (test state))



