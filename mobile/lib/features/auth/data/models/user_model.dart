class UserModel {
  final String id;
  final String phone;
  final String? email;
  final String nickname;
  final String avatar;
  final int age;
  final String gender;
  final String? bio;
  final List<String>? interests;
  final List<String>? photos;
  final GeoLocation? location;
  final VerificationStatus? verification;
  final UserPreferences? preferences;
  final PrivacySettings? privacy;

  const UserModel({
    required this.id,
    required this.phone,
    this.email,
    required this.nickname,
    required this.avatar,
    required this.age,
    required this.gender,
    this.bio,
    this.interests,
    this.photos,
    this.location,
    this.verification,
    this.preferences,
    this.privacy,
  });

  factory UserModel.fromJson(Map<String, dynamic> json) {
    return UserModel(
      id: json['id'] as String,
      phone: json['phone'] as String,
      email: json['email'] as String?,
      nickname: json['nickname'] as String,
      avatar: json['avatar'] as String? ?? '',
      age: json['age'] as int,
      gender: json['gender'] as String,
      bio: json['bio'] as String?,
      interests: (json['interests'] as List<dynamic>?)?.cast<String>(),
      photos: (json['photos'] as List<dynamic>?)?.cast<String>(),
      location: json['location'] != null ? GeoLocation.fromJson(json['location']) : null,
      verification: json['verification'] != null ? VerificationStatus.fromJson(json['verification']) : null,
      preferences: json['preferences'] != null ? UserPreferences.fromJson(json['preferences']) : null,
      privacy: json['privacy'] != null ? PrivacySettings.fromJson(json['privacy']) : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'phone': phone,
      'email': email,
      'nickname': nickname,
      'avatar': avatar,
      'age': age,
      'gender': gender,
      'bio': bio,
      'interests': interests,
      'photos': photos,
      'location': location?.toJson(),
      'verification': verification?.toJson(),
      'preferences': preferences?.toJson(),
      'privacy': privacy?.toJson(),
    };
  }

  UserModel copyWith({
    String? id,
    String? phone,
    String? email,
    String? nickname,
    String? avatar,
    int? age,
    String? gender,
    String? bio,
    List<String>? interests,
    List<String>? photos,
    GeoLocation? location,
    VerificationStatus? verification,
    UserPreferences? preferences,
    PrivacySettings? privacy,
  }) {
    return UserModel(
      id: id ?? this.id,
      phone: phone ?? this.phone,
      email: email ?? this.email,
      nickname: nickname ?? this.nickname,
      avatar: avatar ?? this.avatar,
      age: age ?? this.age,
      gender: gender ?? this.gender,
      bio: bio ?? this.bio,
      interests: interests ?? this.interests,
      photos: photos ?? this.photos,
      location: location ?? this.location,
      verification: verification ?? this.verification,
      preferences: preferences ?? this.preferences,
      privacy: privacy ?? this.privacy,
    );
  }
}

class GeoLocation {
  final double latitude;
  final double longitude;
  final String address;
  final String city;
  final String country;

  const GeoLocation({
    required this.latitude,
    required this.longitude,
    required this.address,
    required this.city,
    required this.country,
  });

  factory GeoLocation.fromJson(Map<String, dynamic> json) {
    return GeoLocation(
      latitude: (json['latitude'] as num).toDouble(),
      longitude: (json['longitude'] as num).toDouble(),
      address: json['address'] as String,
      city: json['city'] as String,
      country: json['country'] as String,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'latitude': latitude,
      'longitude': longitude,
      'address': address,
      'city': city,
      'country': country,
    };
  }
}

class VerificationStatus {
  final IdentityVerification identity;
  final FaceVerification face;

  const VerificationStatus({
    required this.identity,
    required this.face,
  });

  factory VerificationStatus.fromJson(Map<String, dynamic> json) {
    return VerificationStatus(
      identity: IdentityVerification.fromJson(json['identity']),
      face: FaceVerification.fromJson(json['face']),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'identity': identity.toJson(),
      'face': face.toJson(),
    };
  }
}

class IdentityVerification {
  final bool verified;
  final String? idNumber;
  final String? realName;
  final DateTime? verifiedAt;

  const IdentityVerification({
    required this.verified,
    this.idNumber,
    this.realName,
    this.verifiedAt,
  });

  factory IdentityVerification.fromJson(Map<String, dynamic> json) {
    return IdentityVerification(
      verified: json['verified'] as bool,
      idNumber: json['id_number'] as String?,
      realName: json['real_name'] as String?,
      verifiedAt: json['verified_at'] != null ? DateTime.parse(json['verified_at']) : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'verified': verified,
      'id_number': idNumber,
      'real_name': realName,
      'verified_at': verifiedAt?.toIso8601String(),
    };
  }
}

class FaceVerification {
  final bool verified;
  final double confidence;
  final DateTime? verifiedAt;

  const FaceVerification({
    required this.verified,
    required this.confidence,
    this.verifiedAt,
  });

  factory FaceVerification.fromJson(Map<String, dynamic> json) {
    return FaceVerification(
      verified: json['verified'] as bool,
      confidence: (json['confidence'] as num).toDouble(),
      verifiedAt: json['verified_at'] != null ? DateTime.parse(json['verified_at']) : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'verified': verified,
      'confidence': confidence,
      'verified_at': verifiedAt?.toIso8601String(),
    };
  }
}

class UserPreferences {
  final AgeRange ageRange;
  final String genderFilter;
  final int distanceRange;
  final bool showOnline;
  final bool showDistance;

  const UserPreferences({
    required this.ageRange,
    required this.genderFilter,
    required this.distanceRange,
    required this.showOnline,
    required this.showDistance,
  });

  factory UserPreferences.fromJson(Map<String, dynamic> json) {
    return UserPreferences(
      ageRange: AgeRange.fromJson(json['age_range']),
      genderFilter: json['gender_filter'] as String,
      distanceRange: json['distance_range'] as int,
      showOnline: json['show_online'] as bool,
      showDistance: json['show_distance'] as bool,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'age_range': ageRange.toJson(),
      'gender_filter': genderFilter,
      'distance_range': distanceRange,
      'show_online': showOnline,
      'show_distance': showDistance,
    };
  }
}

class AgeRange {
  final int min;
  final int max;

  const AgeRange({
    required this.min,
    required this.max,
  });

  factory AgeRange.fromJson(Map<String, dynamic> json) {
    return AgeRange(
      min: json['min'] as int,
      max: json['max'] as int,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'min': min,
      'max': max,
    };
  }
}

class PrivacySettings {
  final String profileVisibility;
  final bool locationVisible;
  final bool onlineStatus;
  final bool readReceipts;

  const PrivacySettings({
    required this.profileVisibility,
    required this.locationVisible,
    required this.onlineStatus,
    required this.readReceipts,
  });

  factory PrivacySettings.fromJson(Map<String, dynamic> json) {
    return PrivacySettings(
      profileVisibility: json['profile_visibility'] as String,
      locationVisible: json['location_visible'] as bool,
      onlineStatus: json['online_status'] as bool,
      readReceipts: json['read_receipts'] as bool,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'profile_visibility': profileVisibility,
      'location_visible': locationVisible,
      'online_status': onlineStatus,
      'read_receipts': readReceipts,
    };
  }
}