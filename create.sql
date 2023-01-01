CREATE TABLE history (
                         id INTEGER PRIMARY KEY AUTO_INCREMENT,
                         video_url TEXT,
                         start_time TEXT,
                         end_time TEXT,
                         file_name TEXT,
                         is_done BOOLEAN DEFAULT FALSE
);
