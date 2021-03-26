;; MACRO UTILS
(defmacro! default (fn* [arg def] `(if (<= (count optionals) ~arg) ~def (nth optionals ~arg))))
(defmacro! defunc (fn* [name args func] `(def! ~name (fn* ~args ~func))))


;;UTILS
(defunc between [low val max] (&& (>= (ord val) low) (<= (ord val) max)))


;; CONSTANTS
(defunc +number+ [char] (between 48 char 57))


;; READER MACRO

;; PARSER FUNCTION
(defunc read [& optionals]
      	   (let* [stream (default 0 "")
	   	  eof-err (default 1 true)
		  eof-val (default 2 nil)
		  recur (default 3 false)]

		  [stream eof-err eof-val recur (raise "lol")]
		  ))

;; BASIC PARSER SETUP
(def! source (slurp "./parser/main.cl"))

;(prn (macroexpand (default 1 "")))
(prn (read 1))