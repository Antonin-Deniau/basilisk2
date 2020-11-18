; UTILS
(def! add_ast (fn* [prev next]
	{
		:state (:state next)
		:ast (cons (:ast prev) (:ast next)]
	}))

(def! char (fn* [cha]
	(fn* [state]
		(if (= (count (:data state)) 0)
			{ :data "" :ast nil }
			(let* [c (subs (:data state) 0 1)]
				(if (= c cha)
				  {
					:data (subs (:data state) 1)
					:ast c
				  }
				  state)))))) ; RETURN UNCHANGED

(def! choose (fn* [& args]
	(fn* [state]
		(if (= (count args) 0)
			state ; RETURN UNCHANGED
			(let* [res (apply (first args) state)]
				(if (= (:state res) (:state state))
					((apply choose (rest args)) state) 
					(add_ast state res)))))))

(def! repeat (fn* [typ fnc]
	(fn* [state]
		(let* [res (fnc state)]
			(if (= (:state res) (:state state))
				state ; RETURN UNCHANGED
				(let* [resp ((repeat typ fnc) { :data (:data res) :ast ""})]
				  { 
				  	:data (:data resp)
					:ast (apply typ (cons (:ast res) (:ast resp))) }))))))


(def! sequence_inside (fn* [state i & fncs]
	(if (= (count fncs) 0)
	  { :state state :num i }
	  (let* [res ((first fncs) state)]
	    (if (= res state)
	      { :state res :num i }
	      (apply sequence_inside { :state (res) :ast [] } (+ i 1) (rest fncs)))))))

(def! sequence (fn* [typ & args]
	(fn* [state & orig]
		(let* [o (if (= (count orig) 0) state (first orig))
		       res (apply (first args) state)]
			(if (= (count args) 0)
				state
				(if (= res state)
					o ; RETURN ORIGINAL
					(let* [resp ((apply sequence (rest args)) res o)]
						(if (= resp res)
							o ; RETURN ORIGINAL
							{ 
								:data (:data resp)
								:ast (apply typ (cons (:ast res) (:ast resp))) 
							}))))))))

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

;(def! test (sequence vector
;		     (char "6")
;		     (char "9")
;		     (char "6")))

; ENV
(def! data "6969antonin (1 2 4 \"lol\" nil true)")
(def! state { :data data :ast "" })
(prn state)
(prn (sequence_inside state 0
		     (char "6")
		     (char "9")
		     (char "6")))


