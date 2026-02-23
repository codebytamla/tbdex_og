class Patient {
  final String id;
  final String fullName;
  final DateTime? dateOfBirth;
  final String? sex;
  final String? phone;
  final String? nationalId;
  final String? mrn;
  final String status;
  final String? village;

  const Patient({
    required this.id,
    required this.fullName,
    this.dateOfBirth,
    this.sex,
    this.phone,
    this.nationalId,
    this.mrn,
    this.status = 'active',
    this.village,
  });

  int? get age {
    if (dateOfBirth == null) return null;
    final now = DateTime.now();
    int age = now.year - dateOfBirth!.year;
    if (now.month < dateOfBirth!.month ||
        (now.month == dateOfBirth!.month && now.day < dateOfBirth!.day)) {
      age--;
    }
    return age;
  }

  String get sexLabel {
    switch (sex) {
      case 'male':
        return 'Male';
      case 'female':
        return 'Female';
      default:
        return 'Other';
    }
  }

  String get initials {
    final parts = fullName.split(' ');
    if (parts.length >= 2) {
      return '${parts[0][0]}${parts[1][0]}'.toUpperCase();
    }
    return fullName[0].toUpperCase();
  }
}

class TaskItem {
  final String patientId;
  final String patientName;
  final String village;
  final String time;
  final String status; // pending, in_progress, done

  const TaskItem({
    required this.patientId,
    required this.patientName,
    required this.village,
    required this.time,
    required this.status,
  });
}

class AssessmentRecord {
  final String id;
  final String patientId;
  final String patientName;
  final DateTime date;
  final String riskLevel; // low, medium, high
  final bool referralCreated;
  final String templateName;
  final int symptomCount;

  const AssessmentRecord({
    required this.id,
    required this.patientId,
    required this.patientName,
    required this.date,
    required this.riskLevel,
    required this.referralCreated,
    required this.templateName,
    required this.symptomCount,
  });
}

class SymptomItem {
  final String id;
  final String label;
  final String icon;

  const SymptomItem({
    required this.id,
    required this.label,
    required this.icon,
  });
}

class SyncEvent {
  final String message;
  final String time;
  final bool success;
  final bool retryAvailable;

  const SyncEvent({
    required this.message,
    required this.time,
    required this.success,
    this.retryAvailable = false,
  });
}

class AssessmentData {
  double? temperature;
  double? weight;
  int? heartRate;
  int? spo2;
  List<String> selectedSymptomIds;
  bool? tbContactHistory;
  bool? hivPositive;
  bool? previousTbTreatment;
  String notes;

  AssessmentData({
    this.temperature,
    this.weight,
    this.heartRate,
    this.spo2,
    List<String>? selectedSymptomIds,
    this.tbContactHistory,
    this.hivPositive,
    this.previousTbTreatment,
    this.notes = '',
  }) : selectedSymptomIds = selectedSymptomIds ?? [];

  String get riskLevel {
    int score = selectedSymptomIds.length;
    if (tbContactHistory == true) score += 1;
    if (hivPositive == true) score += 2;
    if (previousTbTreatment == true) score += 1;
    // Blood in sputum is always high risk
    if (selectedSymptomIds.contains('s7')) return 'high';
    if (score >= 3) return 'high';
    if (score >= 2) return 'medium';
    return 'low';
  }

  bool get shouldRefer => riskLevel == 'high';

  List<String> get matchedCriteria {
    final criteria = <String>[];
    if (selectedSymptomIds.contains('s1')) criteria.add('Persistent cough');
    if (selectedSymptomIds.contains('s2')) criteria.add('Fever');
    if (selectedSymptomIds.contains('s3')) criteria.add('Night sweats');
    if (selectedSymptomIds.contains('s4')) criteria.add('Weight loss');
    if (selectedSymptomIds.contains('s5')) criteria.add('Chest pain');
    if (selectedSymptomIds.contains('s6')) criteria.add('Fatigue');
    if (selectedSymptomIds.contains('s7')) criteria.add('Blood in sputum');
    if (tbContactHistory == true) criteria.add('TB contact history');
    if (hivPositive == true) criteria.add('HIV positive');
    return criteria;
  }
}
