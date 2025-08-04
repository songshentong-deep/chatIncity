-- 兴趣标签初始化数据
-- 清空现有数据
TRUNCATE TABLE interest_tags;

-- 运动健身类
INSERT INTO interest_tags (id, name, category, description, icon, color, usage_count, is_active, created_at, updated_at) VALUES
('sport-001', '跑步', 'sports', '享受跑步带来的自由与健康', '🏃‍♂️', '#FF6B6B', 0, true, NOW(), NOW()),
('sport-002', '健身', 'sports', '塑造完美身材，保持健康体魄', '💪', '#4ECDC4', 0, true, NOW(), NOW()),
('sport-003', '游泳', 'sports', '在水中感受自由与力量', '🏊‍♀️', '#45B7D1', 0, true, NOW(), NOW()),
('sport-004', '瑜伽', 'sports', '身心合一，寻找内在平衡', '🧘‍♀️', '#96CEB4', 0, true, NOW(), NOW()),
('sport-005', '篮球', 'sports', '团队合作，挥洒汗水', '🏀', '#FFEAA7', 0, true, NOW(), NOW()),
('sport-006', '足球', 'sports', '绿茵场上的激情与梦想', '⚽', '#74B9FF', 0, true, NOW(), NOW()),
('sport-007', '网球', 'sports', '优雅运动，展现风采', '🎾', '#A29BFE', 0, true, NOW(), NOW()),
('sport-008', '羽毛球', 'sports', '轻盈灵动，快乐运动', '🏸', '#FD79A8', 0, true, NOW(), NOW());

-- 音乐类
INSERT INTO interest_tags (id, name, category, description, icon, color, usage_count, is_active, created_at, updated_at) VALUES
('music-001', '流行音乐', 'music', '跟随时代节拍，感受流行魅力', '🎵', '#FF7675', 0, true, NOW(), NOW()),
('music-002', '古典音乐', 'music', '聆听经典，品味高雅艺术', '🎼', '#6C5CE7', 0, true, NOW(), NOW()),
('music-003', '摇滚音乐', 'music', '释放激情，摇摆青春', '🎸', '#2D3436', 0, true, NOW(), NOW()),
('music-004', '民谣', 'music', '简单旋律，诉说生活故事', '🎤', '#00B894', 0, true, NOW(), NOW()),
('music-005', '电子音乐', 'music', '科技与艺术的完美融合', '🎧', '#00CEC9', 0, true, NOW(), NOW()),
('music-006', '爵士乐', 'music', '自由即兴，优雅浪漫', '🎺', '#FDCB6E', 0, true, NOW(), NOW()),
('music-007', '钢琴', 'music', '黑白键上的音乐诗篇', '🎹', '#E17055', 0, true, NOW(), NOW()),
('music-008', 'K歌', 'music', '释放歌声，享受快乐时光', '🎤', '#F39C12', 0, true, NOW(), NOW());

-- 电影类
INSERT INTO interest_tags (id, name, category, description, icon, color, usage_count, is_active, created_at, updated_at) VALUES
('movie-001', '动作片', 'movies', '肾上腺素飙升的视觉盛宴', '🎬', '#E74C3C', 0, true, NOW(), NOW()),
('movie-002', '喜剧片', 'movies', '欢声笑语，轻松愉快', '😂', '#F39C12', 0, true, NOW(), NOW()),
('movie-003', '爱情片', 'movies', '浪漫情怀，温暖人心', '💕', '#E91E63', 0, true, NOW(), NOW()),
('movie-004', '科幻片', 'movies', '探索未来，想象无限', '🚀', '#3498DB', 0, true, NOW(), NOW()),
('movie-005', '悬疑片', 'movies', '烧脑推理，扣人心弦', '🔍', '#9B59B6', 0, true, NOW(), NOW()),
('movie-006', '恐怖片', 'movies', '刺激惊悚，挑战胆量', '👻', '#2C3E50', 0, true, NOW(), NOW()),
('movie-007', '纪录片', 'movies', '真实记录，深度思考', '📹', '#27AE60', 0, true, NOW(), NOW()),
('movie-008', '动画片', 'movies', '童心未泯，想象飞翔', '🎨', '#FF6B6B', 0, true, NOW(), NOW());

-- 阅读类
INSERT INTO interest_tags (id, name, category, description, icon, color, usage_count, is_active, created_at, updated_at) VALUES
('read-001', '小说', 'reading', '沉浸故事世界，体验不同人生', '📚', '#8E44AD', 0, true, NOW(), NOW()),
('read-002', '散文', 'reading', '优美文字，抒发情感', '✍️', '#16A085', 0, true, NOW(), NOW()),
('read-003', '诗歌', 'reading', '韵律之美，意境深远', '📝', '#D35400', 0, true, NOW(), NOW()),
('read-004', '历史', 'reading', '以史为镜，了解过往', '📜', '#7F8C8D', 0, true, NOW(), NOW()),
('read-005', '哲学', 'reading', '思辨人生，探索真理', '🤔', '#34495E', 0, true, NOW(), NOW()),
('read-006', '科普', 'reading', '科学知识，开拓视野', '🔬', '#3498DB', 0, true, NOW(), NOW()),
('read-007', '传记', 'reading', '名人故事，励志人生', '👤', '#E67E22', 0, true, NOW(), NOW()),
('read-008', '心理学', 'reading', '了解内心，认识自我', '🧠', '#9B59B6', 0, true, NOW(), NOW());

-- 旅行类
INSERT INTO interest_tags (id, name, category, description, icon, color, usage_count, is_active, created_at, updated_at) VALUES
('travel-001', '国内游', 'travel', '探索祖国大好河山', '🏔️', '#27AE60', 0, true, NOW(), NOW()),
('travel-002', '出国游', 'travel', '体验异国风情文化', '✈️', '#3498DB', 0, true, NOW(), NOW()),
('travel-003', '自驾游', 'travel', '自由驰骋，说走就走', '🚗', '#E74C3C', 0, true, NOW(), NOW()),
('travel-004', '徒步', 'travel', '用脚步丈量世界', '🥾', '#8E44AD', 0, true, NOW(), NOW()),
('travel-005', '摄影旅行', 'travel', '用镜头记录美好', '📷', '#F39C12', 0, true, NOW(), NOW()),
('travel-006', '海岛游', 'travel', '碧海蓝天，度假天堂', '🏝️', '#00BCD4', 0, true, NOW(), NOW()),
('travel-007', '古镇游', 'travel', '寻找历史足迹，感受古韵', '🏘️', '#795548', 0, true, NOW(), NOW()),
('travel-008', '美食之旅', 'travel', '品尝各地特色美食', '🍜', '#FF5722', 0, true, NOW(), NOW());

-- 美食类
INSERT INTO interest_tags (id, name, category, description, icon, color, usage_count, is_active, created_at, updated_at) VALUES
('food-001', '川菜', 'food', '麻辣鲜香，回味无穷', '🌶️', '#E74C3C', 0, true, NOW(), NOW()),
('food-002', '粤菜', 'food', '清淡鲜美，营养丰富', '🦐', '#27AE60', 0, true, NOW(), NOW()),
('food-003', '日料', 'food', '精致美味，健康养生', '🍣', '#FF9800', 0, true, NOW(), NOW()),
('food-004', '西餐', 'food', '优雅用餐，品味生活', '🥩', '#795548', 0, true, NOW(), NOW()),
('food-005', '烘焙', 'food', '甜蜜制作，分享快乐', '🧁', '#E91E63', 0, true, NOW(), NOW()),
('food-006', '火锅', 'food', '热气腾腾，聚会首选', '🍲', '#F44336', 0, true, NOW(), NOW()),
('food-007', '咖啡', 'food', '香醇浓郁，提神醒脑', '☕', '#8D6E63', 0, true, NOW(), NOW()),
('food-008', '茶艺', 'food', '品茶论道，修身养性', '🍵', '#4CAF50', 0, true, NOW(), NOW());

-- 游戏类
INSERT INTO interest_tags (id, name, category, description, icon, color, usage_count, is_active, created_at, updated_at) VALUES
('game-001', '手机游戏', 'gaming', '随时随地，轻松娱乐', '📱', '#FF9800', 0, true, NOW(), NOW()),
('game-002', '电脑游戏', 'gaming', '沉浸体验，畅快游戏', '💻', '#2196F3', 0, true, NOW(), NOW()),
('game-003', '主机游戏', 'gaming', '专业设备，极致体验', '🎮', '#9C27B0', 0, true, NOW(), NOW()),
('game-004', '桌游', 'gaming', '面对面交流，增进友谊', '🎲', '#FF5722', 0, true, NOW(), NOW()),
('game-005', '棋牌', 'gaming', '智力对决，策略思考', '♠️', '#607D8B', 0, true, NOW(), NOW()),
('game-006', 'VR游戏', 'gaming', '虚拟现实，未来体验', '🥽', '#00BCD4', 0, true, NOW(), NOW()),
('game-007', '竞技游戏', 'gaming', '团队协作，竞技对抗', '🏆', '#FFC107', 0, true, NOW(), NOW()),
('game-008', '休闲游戏', 'gaming', '轻松愉快，放松心情', '🎯', '#4CAF50', 0, true, NOW(), NOW());

-- 艺术类
INSERT INTO interest_tags (id, name, category, description, icon, color, usage_count, is_active, created_at, updated_at) VALUES
('art-001', '绘画', 'art', '用画笔描绘内心世界', '🎨', '#E91E63', 0, true, NOW(), NOW()),
('art-002', '书法', 'art', '笔墨纸砚，传承文化', '✒️', '#424242', 0, true, NOW(), NOW()),
('art-003', '雕塑', 'art', '立体艺术，空间美学', '🗿', '#8D6E63', 0, true, NOW(), NOW()),
('art-004', '摄影', 'art', '定格美好，记录瞬间', '📸', '#FF9800', 0, true, NOW(), NOW()),
('art-005', '舞蹈', 'art', '身体语言，优美表达', '💃', '#E91E63', 0, true, NOW(), NOW()),
('art-006', '戏剧', 'art', '舞台表演，情感演绎', '🎭', '#9C27B0', 0, true, NOW(), NOW()),
('art-007', '手工', 'art', '亲手制作，创意无限', '✂️', '#FF5722', 0, true, NOW(), NOW()),
('art-008', '设计', 'art', '美学创造，实用艺术', '📐', '#2196F3', 0, true, NOW(), NOW());

-- 科技类
INSERT INTO interest_tags (id, name, category, description, icon, color, usage_count, is_active, created_at, updated_at) VALUES
('tech-001', '编程', 'technology', '代码世界，逻辑思维', '💻', '#4CAF50', 0, true, NOW(), NOW()),
('tech-002', '人工智能', 'technology', '智能未来，科技前沿', '🤖', '#2196F3', 0, true, NOW(), NOW()),
('tech-003', '数码产品', 'technology', '科技生活，便捷体验', '📱', '#FF9800', 0, true, NOW(), NOW()),
('tech-004', '区块链', 'technology', '去中心化，技术革命', '⛓️', '#9C27B0', 0, true, NOW(), NOW()),
('tech-005', '物联网', 'technology', '万物互联，智慧生活', '🌐', '#00BCD4', 0, true, NOW(), NOW()),
('tech-006', '云计算', 'technology', '云端服务，无限可能', '☁️', '#607D8B', 0, true, NOW(), NOW()),
('tech-007', '大数据', 'technology', '数据分析，洞察未来', '📊', '#FF5722', 0, true, NOW(), NOW()),
('tech-008', '网络安全', 'technology', '信息保护，安全第一', '🔒', '#F44336', 0, true, NOW(), NOW());

-- 时尚类
INSERT INTO interest_tags (id, name, category, description, icon, color, usage_count, is_active, created_at, updated_at) VALUES
('fashion-001', '服装搭配', 'fashion', '时尚穿搭，展现个性', '👗', '#E91E63', 0, true, NOW(), NOW()),
('fashion-002', '美妆', 'fashion', '精致妆容，美丽自信', '💄', '#FF5722', 0, true, NOW(), NOW()),
('fashion-003', '护肤', 'fashion', '肌肤保养，由内而外', '🧴', '#4CAF50', 0, true, NOW(), NOW()),
('fashion-004', '珠宝', 'fashion', '璀璨饰品，点缀生活', '💎', '#9C27B0', 0, true, NOW(), NOW()),
('fashion-005', '包包', 'fashion', '时尚配件，实用美观', '👜', '#795548', 0, true, NOW(), NOW()),
('fashion-006', '鞋子', 'fashion', '足下风采，舒适时尚', '👠', '#424242', 0, true, NOW(), NOW()),
('fashion-007', '香水', 'fashion', '迷人香气，个性标签', '🌸', '#E91E63', 0, true, NOW(), NOW()),
('fashion-008', '潮流', 'fashion', '紧跟时尚，引领潮流', '✨', '#FF9800', 0, true, NOW(), NOW());

-- 宠物类
INSERT INTO interest_tags (id, name, category, description, icon, color, usage_count, is_active, created_at, updated_at) VALUES
('pet-001', '猫咪', 'pets', '可爱猫咪，温暖陪伴', '🐱', '#FF9800', 0, true, NOW(), NOW()),
('pet-002', '狗狗', 'pets', '忠诚伙伴，快乐源泉', '🐶', '#8D6E63', 0, true, NOW(), NOW()),
('pet-003', '鸟类', 'pets', '自由飞翔，歌声悦耳', '🐦', '#4CAF50', 0, true, NOW(), NOW()),
('pet-004', '鱼类', 'pets', '静谧观赏，心灵宁静', '🐠', '#00BCD4', 0, true, NOW(), NOW()),
('pet-005', '仓鼠', 'pets', '小巧可爱，活泼机灵', '🐹', '#FF5722', 0, true, NOW(), NOW()),
('pet-006', '兔子', 'pets', '温顺可爱，毛茸茸的', '🐰', '#E91E63', 0, true, NOW(), NOW()),
('pet-007', '爬虫', 'pets', '独特魅力，神秘伙伴', '🦎', '#4CAF50', 0, true, NOW(), NOW()),
('pet-008', '宠物护理', 'pets', '用心照顾，健康成长', '🏥', '#2196F3', 0, true, NOW(), NOW());

-- 摄影类
INSERT INTO interest_tags (id, name, category, description, icon, color, usage_count, is_active, created_at, updated_at) VALUES
('photo-001', '人像摄影', 'photography', '捕捉美丽，定格瞬间', '📷', '#E91E63', 0, true, NOW(), NOW()),
('photo-002', '风景摄影', 'photography', '大自然美景，壮丽画面', '🏞️', '#4CAF50', 0, true, NOW(), NOW()),
('photo-003', '街拍', 'photography', '城市生活，真实记录', '🏙️', '#607D8B', 0, true, NOW(), NOW()),
('photo-004', '微距摄影', 'photography', '细节之美，放大世界', '🔍', '#FF9800', 0, true, NOW(), NOW()),
('photo-005', '黑白摄影', 'photography', '经典艺术，情感表达', '⚫', '#424242', 0, true, NOW(), NOW()),
('photo-006', '夜景摄影', 'photography', '夜色朦胧，光影魅力', '🌃', '#2196F3', 0, true, NOW(), NOW()),
('photo-007', '婚礼摄影', 'photography', '幸福时刻，永恒记忆', '💒', '#E91E63', 0, true, NOW(), NOW()),
('photo-008', '后期制作', 'photography', '数字艺术，完美呈现', '🎨', '#9C27B0', 0, true, NOW(), NOW());