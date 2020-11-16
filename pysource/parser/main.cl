; UTILS
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
				  state)))))) ; RETURN UNCHANGED

(def! choose (fn* [& args]
	(fn* [state]
		(if (= (count args) 0)
			state ; RETURN UNCHANGED
			(let* [res (apply (first args) state)]
				(if (= res state)
					((apply choose (rest args)) state) res))))))

(def! repeat (fn* [typ fnc]
	(fn* [state]
		(let* [res (fnc state)]
			(if (= res state)
				state ; RETURN UNCHANGED
				(let* [resp ((repeat typ fnc) { :data (:data res) :ast ""})]
				  { 
				  	:data (:data resp)
					:ast (apply typ (cons (:ast res) (:ast resp))) }))))))

(def! sequence (fn* [& args]
	(fn* [state & orig]
		(let* [res (apply (first args) state)]
			(if (= (count args) 1)
				res
				(if (= res state)
					(if (= orig ()) state (first orig)) ; RETURN ORIGINAL
					((apply sequence (rest args)) res (if (= orig ()) state (first orig)))
					))))))

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
				state))))) ; RETURN UNCHANGED

(def! ALPHA (fn* [state]
	(if (= (count (:data state)) 0)
		{ :data "" :ast (:ast state) }

		(let* [c (ord (subs (:data state) 0 1))]
			(if (&& (<= 65 c) (<= c 90))
				(nums {
					:data (subs (:data state) 1)
					:ast (str (:ast state) (chr c))
				})
				state))))) ; RETURN UNCHANGED

(def! alpha (fn* [state]
	(if (= (count (:data state)) 0)
		{ :data "" :ast (:ast state) }

		(let* [c (ord (subs (:data state) 0 1))]
			(if (&& (<= 97 c) (<= c 122))
				(nums {
					:data (subs (:data state) 1)
					:ast (str (:ast state) (chr c))
				})
				state))))) ; RETURN UNCHANGED

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

(def! test (sequence (char "6")
		     (char "9")
		     (char "6")))

; ENV
(def! data "6969antonin (1 2 4 \"lol\" nil true)")
(def! state { :data data :ast "" })
(prn state)
(prn (test state))



