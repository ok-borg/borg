select id, title from posts where tags like '%<bash>%' AND score >= 4 order by score DESC;
