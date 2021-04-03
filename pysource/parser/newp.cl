;; MACRO UTILS
(defmacro! default (fn* [arg def] `(if (nil? ~arg) ~def ~arg)))
(defmacro! defunc (fn* [name args func] `(def! ~name (fn* ~args ~func))))

;;UTILS
(defunc ord-between [low val max] (&& (>= (ord val) low) (<= (ord val) max)))


;; MATCHER
(defunc number-matcher [c] (ord-between 48 c 57))
(defunc vector-matcher [c] (= "[" c))
(defunc whitespace-matcher [c] (reduce || [(= " " c) (= "\n" c) (= "\t" c)]))


;; READER
(defunc whitespace-ignore [stream] 
        (if (whitespace-matcher (peek-byte stream))
          (do
            (read-byte stream)
            (whitespace-ignore stream))
          stream))

(defunc number-reader [reader-macro stream] (read-byte stream))

(defunc vector-reader-iterate [ret reader-macro stream]
        (let* [stream (whitespace-ignore stream)] ;; IGNORE WHITESPACE
          (if (= (peek-byte stream) "]") ;; RETURN ON VECTOR END
            ret
            (let* [data (read reader-macro stream)]
                  (vector-reader-iterate (conj ret data) reader-macro stream)))))

(defunc vector-reader [reader-macro stream] 
        (do
          (read-byte stream)
          (vector-reader-iterate [] reader-macro stream)))


;; READER MACRO
(def! reader-macro [
    [number-matcher number-reader]
    [vector-matcher vector-reader]])

;; PARSER FUNCTION
(defunc match [reader-macro c]
          (if (empty? reader-macro)
            nil
            (let* [macro (first reader-macro)
                         matcher (nth macro 0)
                         reader  (nth macro 1)
                         is-true (matcher c)]
              (if is-true
                reader
                (match (rest reader-macro) c)))))

(defunc read [reader-macro stream]
        (let* [stream (whitespace-ignore stream)
               byte   (peek-byte stream)
               reader (match reader-macro byte)]
          (if (nil? reader)
            (raise (str "Unable to find matcher: [" byte "]"))
            (reader reader-macro stream))))

;; BASIC PARSER SETUP
(prn (read reader-macro (string-stream " [ -1.5554 .2 5484.263 5 ] ")))
