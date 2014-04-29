//
//  WHOMessageTableViewController.h
//  Hooli
//
//  Created by dylan on 4/19/14.
//  Copyright (c) 2014 whoisdylan. All rights reserved.
//

#import <UIKit/UIKit.h>
#import <CoreLocation/CoreLocation.h>

@interface WHOMessageTableViewController : UITableViewController <CLLocationManagerDelegate>
- (id)initWithStyle:(UITableViewStyle)style WithUserName:(NSString* )userName WithEncodedPhoto:(NSString* )encodedPhoto;

@property (nonatomic, strong) NSMutableArray* messages;
@property (nonatomic, strong) NSString* userName;
@property (nonatomic, strong) CLLocationManager *locationManager;
@property (nonatomic, strong) CLLocation* userLocation;
@property (nonatomic, strong) NSString* userEncodedPhoto;

@end
