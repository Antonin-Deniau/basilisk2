
(defmacro! each (fn* [args body]
   `(let* [fnc (fn* [& args]
                    (do
                     (let* [(first args) (first args)] ~body)
                     (if (&& (empty? args))
                       nil
                       (apply fnc (rest args)))))]
      (fnc ~@(nth args 1))
     )))

;(prn (macroexpand

(each [i [1 2 3 4]]
      (prn i))

;))
