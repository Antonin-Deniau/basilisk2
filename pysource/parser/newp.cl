;; MACRO UTILS
(defmacro! default (fn* [arg def] `(if (nil? ~arg) ~def ~arg)))
(defmacro! defunc (fn* [name args func] `(def! ~name (fn* ~args ~func))))

;;UTILS
(defunc between [low val max] (&& (>= (ord val) low) (<= (ord val) max)))


;; MATCHER
(defunc number-matcher [char] (between 48 char 57))
(defunc vector-matcher [] 1)


;; READER
(defunc number-reader [stream] 1)
(defunc vector-reader [stream] 1)

;; READER MACRO
(def! reader-macro [
    [number-matcher number-reader]
    [vector-matcher vector-reader]
  ])

;; PARSER FUNCTION
(defunc read [reader-macro stream]
        (if (empty? reader-macro)
          (raise "Unable to find matcher")
          (let* [macro (first reader-macro)
                       char (peek-byte stream)
                       is-true ((nth macro 0) char)]
            (if is-true
              ((nth macro 1) stream)
              (read (rest reader-macro) stream)))))

;; BASIC PARSER SETUP
(prn (read reader-macro (string-stream " [ 1 2 3 ] ")))
