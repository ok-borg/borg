select body, score, ParentId from posts where ParentId IN (select id from posts where tags like '%<bash>%' AND score >= 4);
