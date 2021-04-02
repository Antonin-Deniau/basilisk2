;; MACRO UTILS
(defmacro! default (fn* [arg def] `(if (nil? ~arg) ~def ~arg)))
(defmacro! defunc (fn* [name args func] `(def! ~name (fn* ~args ~func))))

;;UTILS
(defunc between [low val max] (&& (>= (ord val) low) (<= (ord val) max)))


;; MATCHER
(defunc number-matcher [char] (between 48 char 57))
(defunc vector-matcher [] 1)


;; READER
(defunc number-reader [] 1)
(defunc vector-reader [] 1)

;; READER MACRO
(def! reader-macro [
    [number-matcher number-reader]
    [vector-matcher vector-reader]
  ])

;; PARSER FUNCTION
(defunc read [& optionals]
        (let* [stream  (default (nth opt 0 ""))
               eof-err (default (nth opt 1 true))
               eof-val (default (nth opt 2 nil))
               recur   (default (nth opt 3 false))]

          [stream eof-err eof-val recur "lol"]
          ))

;; BASIC PARSER SETUP
(def! source (slurp "./parser/main.cl"))

;(prn (macroexpand (default 1 "")))
;(prn (read 1))
