;; MACRO UTILS
(defmacro! default (fn* [arg def] `(if (nil? ~arg) ~def ~arg)))
(defmacro! defunc (fn* [name args func] `(def! ~name (fn* ~args ~func))))

;;UTILS
(defunc between [low val max] (&& (>= (ord val) low) (<= (ord val) max)))


;; MATCHER
(defunc number-matcher [c] (between 48 c 57))
(defunc vector-matcher [c] (= "[" c))
(defunc whitespace-matcher [c] (|| (|| (= " " c) ( "\n" c)) (= "\t" c)))


;; READER
(defunc whitespace-reader [reader-macro stream] 
        (if (whitespace-matcher (peek-byte stream))
          (do
            (read-byte stream)
            (whitespace-reader reader-macro stream))
          stream))

(defunc number-reader [reader-macro stream] 2)


(defunc vector-reader-iterate [ret reader-macro stream]
        (let* [stream (whitespace-reader reader-macro stream)] ;; IGNORE WHITESPACE
          (if (= (peek-byte stream) "]") ;; RETURN ON VECTOR END
            ret
            (vector-reader-iterate (conj ret (read reader-macro stream)) reader-macro stream))))

(defunc vector-reader [reader-macro stream] 
        (do
          (read-byte stream) ;; IGNORE [
          (vector-reader-iterate [] reader-macro stream)))


;; READER MACRO
(def! reader-macro [
    [number-matcher number-reader]
    [vector-matcher vector-reader]
  ])

;; PARSER FUNCTION
(defunc read [reader-macro stream]
        (let* [stream (whitespace-reader reader-macro stream)]
          (if (empty? reader-macro)
            (raise "Unable to find matcher")
            (let* [macro   (first reader-macro)
                           matcher (nth macro 0)
                           reader  (nth macro 1)
                           char    (peek-byte stream)
                           is-true (matcher char)]
              (if is-true
                (reader reader-macro stream)
                (read (rest reader-macro) stream))))))

;; BASIC PARSER SETUP
(prn (read reader-macro (string-stream " [ 1 2 3 ] ")))
