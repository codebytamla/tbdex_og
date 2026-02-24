import '../models/models.dart';

class MockData {
  static final List<Patient> patients = [
    Patient(
      id: '1',
      fullName: 'Grace Achieng',
      dateOfBirth: DateTime(1992, 5, 14),
      sex: 'female',
      phone: '+254 712 345 678',
      nationalId: '29384756',
      mrn: 'KER-0042',
      status: 'active',
      village: 'Kericho',
    ),
    Patient(
      id: '2',
      fullName: 'John Mwangi',
      dateOfBirth: DateTime(1985, 11, 3),
      sex: 'male',
      phone: '+254 723 456 789',
      nationalId: '18273645',
      mrn: 'KIP-0018',
      status: 'active',
      village: 'Kipkelion',
    ),
    Patient(
      id: '3',
      fullName: 'Esther Kimani',
      dateOfBirth: DateTime(1998, 2, 20),
      sex: 'female',
      phone: '+254 734 567 890',
      nationalId: '37485926',
      mrn: 'ELD-0071',
      status: 'active',
      village: 'Eldoret',
    ),
    Patient(
      id: '4',
      fullName: 'David Omondi',
      dateOfBirth: DateTime(1970, 8, 12),
      sex: 'male',
      phone: '+254 745 678 901',
      nationalId: '46375819',
      mrn: 'NAK-0033',
      status: 'active',
      village: 'Nakuru',
    ),
    Patient(
      id: '5',
      fullName: 'Mary Wanjiku',
      dateOfBirth: DateTime(2001, 1, 5),
      sex: 'female',
      phone: '+254 756 789 012',
      nationalId: '55748392',
      mrn: 'KER-0056',
      status: 'provisional',
      village: 'Kericho',
    ),
    Patient(
      id: '6',
      fullName: 'James Otieno',
      dateOfBirth: DateTime(1963, 7, 22),
      sex: 'male',
      phone: '+254 767 890 123',
      nationalId: '12345678',
      mrn: 'KIP-0025',
      status: 'active',
      village: 'Kipkelion',
    ),
  ];

  static final List<TaskItem> todaysTasks = [
    TaskItem(
      patientId: '2',
      patientName: 'John Mwangi',
      village: 'Kipkelion',
      time: '9:00 AM',
      status: 'pending',
    ),
    TaskItem(
      patientId: '3',
      patientName: 'Esther Kimani',
      village: 'Eldoret',
      time: '10:30 AM',
      status: 'in_progress',
    ),
    TaskItem(
      patientId: '4',
      patientName: 'David Omondi',
      village: 'Nakuru',
      time: '1:00 PM',
      status: 'done',
    ),
    TaskItem(
      patientId: '1',
      patientName: 'Grace Achieng',
      village: 'Kericho',
      time: '2:30 PM',
      status: 'pending',
    ),
    TaskItem(
      patientId: '6',
      patientName: 'James Otieno',
      village: 'Kipkelion',
      time: '3:45 PM',
      status: 'pending',
    ),
  ];

  static final List<AssessmentRecord> assessmentHistory = [
    AssessmentRecord(
      id: 'a1',
      patientId: '4',
      patientName: 'David Omondi',
      date: DateTime(2026, 2, 23, 13, 0),
      riskLevel: 'low',
      referralCreated: false,
      templateName: 'Standard TB Screening v2.1',
      symptomCount: 1,
    ),
    AssessmentRecord(
      id: 'a2',
      patientId: '1',
      patientName: 'Grace Achieng',
      date: DateTime(2026, 2, 22, 10, 30),
      riskLevel: 'high',
      referralCreated: true,
      templateName: 'Standard TB Screening v2.1',
      symptomCount: 4,
    ),
    AssessmentRecord(
      id: 'a3',
      patientId: '3',
      patientName: 'Esther Kimani',
      date: DateTime(2026, 2, 21, 14, 15),
      riskLevel: 'medium',
      referralCreated: false,
      templateName: 'Standard TB Screening v2.1',
      symptomCount: 2,
    ),
    AssessmentRecord(
      id: 'a4',
      patientId: '2',
      patientName: 'John Mwangi',
      date: DateTime(2026, 2, 20, 9, 0),
      riskLevel: 'low',
      referralCreated: false,
      templateName: 'Pediatric TB Screening v1.0',
      symptomCount: 0,
    ),
  ];

  static final List<SymptomItem> tbSymptoms = [
    SymptomItem(id: 's1', label: 'Persistent cough (>2 weeks)', icon: 'lungs'),
    SymptomItem(id: 's2', label: 'Fever', icon: 'thermometer'),
    SymptomItem(id: 's3', label: 'Night sweats', icon: 'moon'),
    SymptomItem(id: 's4', label: 'Weight loss', icon: 'trending_down'),
    SymptomItem(id: 's5', label: 'Chest pain', icon: 'heart'),
    SymptomItem(id: 's6', label: 'Fatigue', icon: 'battery_low'),
    SymptomItem(id: 's7', label: 'Blood in sputum', icon: 'droplet'),
  ];

  static final SyncStatus syncStatus = SyncStatus(
    isOnline: true,
    pendingUploadCount: 0,
    pendingDownloadCount: 0,
    failedItemCount: 0,
    lastSyncedAt: DateTime.now().subtract(const Duration(minutes: 18)),
  );

  static final List<SyncEvent> syncHistory = [
    SyncEvent(
      message: '5 assessments uploaded',
      time: '10 min ago',
      success: true,
    ),
    SyncEvent(message: 'Templates updated', time: '2 hours ago', success: true),
    SyncEvent(
      message: '1 referral uploaded',
      time: '3 hours ago',
      success: true,
    ),
  ];
}
