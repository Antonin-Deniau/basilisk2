; UTILS
(def! getaddr (fn* [data index]
		   { :line 1 :col 2 }
))

(def! chr (fn* [state]

))

(def! choose (fn* [args]
	(fn* [state]
		(if (= (count args) 0)
			{ :valid false :message (prn "Could not find matching rules in: " args)}

			(let* [res ((get (first args)) state)]
				(if (true? (get res :valid))
				  res
				  (choose (rest args) state)
				)
			)
		)
	)
))

(def! repeat (fn* []
		 1
		 ))
; SYNTAX

(def! nums (fn* [state]
	(if (= (count (get state :data)) 0)
		; TODO { :valid false :message (prn "Reached EOF " args)}

		(let* [char (ord (first (get state :data)))]
			(if (&& (> char 32) (< char 127))
			  (alpha (rest args) state)
			  res)))))

(def! alpha (fn* [state]
))

(def! whitespace (fn* [state]
))

(def! blist (fn* []
	(sequence (chr "(")
		  (repeat bexpr)
		  (chr ")"))))

(def! bkeyword (fn* []
	(sequence (chr ":")
		  (choose alpha
			  (chr "_") 
			  (chr "-"))
		  (repeat (choose alpha
				  nums
				  (chr "_") 
				  (chr "-"))))))

(def! bexpr (fn* []
	(choose blist
		bkeyword)))

; ENV
(def! data "(1 2 4 \"lol\" nil true)")
(def! state { :data data :ast nil })

(prn (bexpr state))

