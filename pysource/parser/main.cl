;(defmacro! defun! (fn* [name args body] `(def! ~name (fn* ~args ~body))))

; UTILS
(def! add_ast (fn* [typ prev next]
                   {
                   :data (:data next)
                   :ast (if (= (:ast prev nil) nil)
                          (:ast next)
                          (apply typ (:ast prev) (:ast next)))
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

(defmacro! choose (fn* [& args]
                  `(fn* [state]
                       (if (= (count ~args) 0)
                         state ; RETURN UNCHANGED
                         (let* [res (apply (first ~args) state)]
                           (if (= (:data res) (:data state))
                             ((apply choose (rest ~args)) state) 
                             res))))))

(defmacro! repeat (fn* [typ fnc]
                  `(fn* [state]
                       (let* [res (~fnc state)]
                         (if (= (:data res) (:data state))
                           state ; RETURN UNCHANGED
                           (let* [resp ((repeat ~typ ~fnc) { :data (:data res) :ast ""})]
                             { 
                             :data (:data resp)
                             :ast (apply ~typ (cons (:ast res) (:ast resp))) }))))))

(def! sequence_inside (fn* [state typ & fncs]
                           (if (= (count fncs) 0)
                             { :data (:data state) :valid true :ast (:ast state) }
                             (let* [res ((first fncs) state)]
                               (if (= (:data res) (:data state))
                                 { :data res :valid false :ast (:ast state) }
                                 (apply sequence_inside (add_ast typ state res) typ (rest fncs)))))))

(defmacro! sequence (fn* [typ & fncs]
                    `(fn* [state & orig]
                         (let* [o (if (= (count orig) 0) state (first orig))
                                  res (apply sequence_inside state ~typ ~fncs)]
                           (if (= (:valid res) true) res o)))))

(defmacro! ignore (fn* [func fncs]
                  `(fn* [state]
                       (let* [res (~fncs state)]
                         (let* [r (~func res)] 
                           {
                           :data (:data r)
                           :ast (:ast res)
                           })
                         ))))

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


(def! whitespace (repeat str
                         (choose (char " ")
                                 (char "	")
                                 (char "\n")
                                 (char ","))))


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


;(def! bnum (fn* [] 
;                1
;                ))
;(de! bnil (fn* []
;               1
;           ))
;(def! blist (sequence (ignore whitespace (char "("))
;                      (repeat list bexpr)
;                      (char "(")))
;(def! bexpr (ignore whitespace 
;                    (choose blist bnum bnil bbool)))
;(def! ast (repeat list bexpr))
; ENV
;(prn (ast { :data data :ast nil }))

;(defun! read [a b & c d] [a b c d])
;(prn (read 1 2 3 4))
;(prn (read 1 2 3 4 5 6))
;(prn (read 1 2))

(defun! parse_list []
        )

(def! macro_chars {
        "(" parse_list
      })

(prn macro_chars)

(defun! read [stream]
        1
        )

(def! data "(1 2 4 nil true)")
(prn (read data))


